package yt_stats

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
)

// Sorts the comments slice based on comment and reply publishing time.
func SortComments(comments *[]interface{}) {
	sortFunc := func(i int, j int) bool {
		switch com1 := (*comments)[i].(type) {
		case Comment:
			switch com2 := (*comments)[j].(type) {
			case Comment:
				return com1.PublishedAt < com2.PublishedAt
			case Reply:
				return com1.PublishedAt < com2.PublishedAt
			}
		case Reply:
			switch com2 := (*comments)[j].(type) {
			case Comment:
				return com1.PublishedAt < com2.PublishedAt
			case Reply:
				return com1.PublishedAt < com2.PublishedAt
			}
		}
		return false
	}
	sort.Slice(*comments, sortFunc)
}

// Searches for one or more substrings in a text. Used for filtering based on content and user.
func searchContent(substrings []string, message string, caseSensitive bool) (bool, bool) {
	any := false
	all := true
	if !caseSensitive {
		message = strings.ToLower(message)
	}
	for _, substring := range substrings {
		if !caseSensitive {
			substring = strings.ToLower(substring)
		}
		if strings.Contains(message, substring) {
			any = true
		} else {
			all = false
		}
	}
	return any, all
}

// The logic used to filter through comments and replies based on multiple searches.
// Returns a bool on if the filtering succeeded, and the result. Cannot check for nil result since empty results exist.
func CommentFilter(searches []Search, comments []interface{}) (bool, []interface{}) {
	if searches == nil {
		return true, comments
	}
	var matches []interface{}
	var source []interface{}
	for _, search := range searches {
		var match []interface{}
		var remains []interface{}
		if search.Reductive {  // Reductive filter, all non-matches stay.
			source = matches
			remains = comments
		} else {  // Additive filter, all matches stay.
			source = comments
			match = matches
		}
		for _, item := range source {
			var msg string
			var name string
			switch com := item.(type) {
			case Comment:
				msg = com.Message
				name = com.AuthorName
			case Reply:
				msg = com.Message
				name = com.AuthorName
			default:
				return false, nil
			}
			anyContent, allContent := searchContent(search.Content, msg, search.CaseSensitive)
			anyUser, allUser := searchContent(search.Users, name, search.CaseSensitive)
			if (search.MatchAny && (anyContent || anyUser)) || allContent && allUser {
				match = append(match, item)
			} else {
				remains = append(remains, item)
			}
		}
		matches = match
		comments = remains
	}
	return true, matches
}

// Worker function that gets replies for comments from a channel of comment IDs. Handles pagination of replies.
// Parses retrieved replies into the comments slice if no errors are found. Otherwise drains channel to save on quota.
// Error or generic OK StatusCodeOutbound struct is deposited into channel to preserve and propagate errors received.
func worker(in <-chan string, c *[]interface{}, r chan<- StatusCodeOutbound, m *sync.Mutex, inp Inputs, k string) int {
	quota := 0
	for comId := range in {
		pageToken := ""
		for hasNextPage := true; hasNextPage; hasNextPage = pageToken != "" {
			var youtubeStatus StatusCodeOutbound
			var repliesInbound RepliesInbound
			ok := func() StatusCodeOutbound { // Function for deferring the closing of response bodies inside loop.
				resp, err := http.Get(fmt.Sprintf("%s&parentId=%s&key=%s&pageToken=%s",
					inp.RepliesRoot, comId, k, pageToken))
				if err != nil {
					return StatusCodeOutbound{
						StatusCode:    http.StatusInternalServerError,
						StatusMessage: "failedToQueryYouTubeAPI",
					}
				}
				defer resp.Body.Close()
				quota += 2  // Snippet cost for this endpoint is 1 less than everywhere else.
				youtubeStatus = ErrorParser(resp.Body, &repliesInbound)
				if youtubeStatus.StatusCode != http.StatusOK {
					return youtubeStatus
				}
				m.Lock()
				RepliesParser(repliesInbound, c)
				m.Unlock()
				pageToken = repliesInbound.NextPageToken
				return youtubeStatus
			}
			if ret := ok(); ret.StatusCode != http.StatusOK {
				r <- ret
				for range in {} // Encountered an error, drain channel to save quota.
				return quota
			}
		}
	}
	r <- StatusCodeOutbound{
		StatusCode:    http.StatusOK,
		StatusMessage: "ok",
	}
	return quota
}

func CommentsHandler(input Inputs) http.Handler {
	workers := 10
	comments := func(w http.ResponseWriter, r *http.Request) {
		quota := 0
		switch r.Method {
		case http.MethodGet:
			key := r.URL.Query().Get("key")
			if key == "" {
				sendStatusCode(w, quota, http.StatusBadRequest, "keyMissing")
				return
			}
			id := r.URL.Query().Get("id")
			if id == "" {
				sendStatusCode(w, quota, http.StatusBadRequest, "videoIdMissing")
				return
			}
			var searches []Search
			if r.Body != nil {
				r.Body = http.MaxBytesReader(w, r.Body, 1048576) // Read max 1 MB
				searchErr := json.NewDecoder(r.Body).Decode(&searches)
				if searchErr != nil && searchErr.Error() == "http: request body too large" {
					sendStatusCode(w, quota, http.StatusRequestEntityTooLarge, "searchBodyTooLarge")
					return
				} else if searchErr != nil && searchErr != io.EOF {
					sendStatusCode(w, quota, http.StatusBadRequest, "searchBodyInvalid")
					return
				}
			}
			var commentsOutbound CommentOutbound
			commentsOutbound.VideoId = id
			var comments []interface{}
			var needReplies []string
			pageToken := ""
			for hasNextPage := true; hasNextPage; hasNextPage = pageToken != "" {
				var youtubeStatus StatusCodeOutbound
				var commentsInbound CommentsInbound
				ok := func() bool { // Internal function for deferring the closing of response bodies inside loop.
					resp, err := http.Get(fmt.Sprintf("%s&videoId=%s&key=%s&pageToken=%s",
						input.CommentsRoot, id, key, pageToken))
					if err != nil {
						sendStatusCode(w, quota, http.StatusInternalServerError, "failedToQueryYouTubeAPI")
						return false
					}
					defer resp.Body.Close()
					quota += 5
					youtubeStatus = ErrorParser(resp.Body, &commentsInbound)
					if youtubeStatus.StatusCode != http.StatusOK {
						if youtubeStatus.StatusMessage == "keyInvalid" {  // Quota cannot be deducted from invalid keys.
							quota -= 5
						}
						sendStatusCode(w, quota, youtubeStatus.StatusCode, youtubeStatus.StatusMessage)
						return false
					}
					CommentsParser(commentsInbound, &comments, &needReplies)
					pageToken = commentsInbound.NextPageToken
					return true
				}
				if !ok() {
					return
				}
			}
			replyIds := make(chan string, len(needReplies))
			for _, comId := range needReplies {
				replyIds <- comId
			}
			close(replyIds)
			var wg sync.WaitGroup
			var mut sync.Mutex
			var add sync.Mutex
			workerResponses := make(chan StatusCodeOutbound, workers)
			wg.Add(workers)
			for i := 0; i < workers; i++ {  // Launch workers.
				go func() {
					n := worker(replyIds, &comments, workerResponses, &mut, input, key)
					add.Lock()
					quota += n
					add.Unlock()
					wg.Done()
				}()
			}
			wg.Wait()
			close(workerResponses)
			for response := range workerResponses {
				if response.StatusCode != http.StatusOK {
					sendStatusCode(w, quota, response.StatusCode, response.StatusMessage)
					return
				}
			}
			ok, filteredComments := CommentFilter(searches, comments)
			if !ok {
				sendStatusCode(w, quota, http.StatusInternalServerError, "failedFilteringComments")
				return
			}
			SortComments(&comments)
			commentsOutbound.Comments = filteredComments
			commentsOutbound.QuotaUsage = quota
			err := json.NewEncoder(w).Encode(commentsOutbound)
			if err != nil {
				log.Println("Failed to respond to playlist endpoint.")
			}
			return
		default:
			unsupportedRequestType(w)
			return
		}
	}
	return http.HandlerFunc(comments)
}
