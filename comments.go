package yt_stats

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
)

// Sorts the comments slice based on comment and reply publishing time.
func SortComments(comments *[]interface{}) {
	sortFunc := func(i int, j int) bool {
		switch com1 := (*comments)[i].(type) {
		case Comment:
			switch com2 := (*comments)[j].(type) {
			case Comment: return com1.PublishedAt < com2.PublishedAt
			case Reply: return com1.PublishedAt < com2.PublishedAt
			}
		case Reply:
			switch com2 := (*comments)[j].(type) {
			case Comment: return com1.PublishedAt < com2.PublishedAt
			case Reply: return com1.PublishedAt < com2.PublishedAt
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
func CommentSearch(searches []Search, comments []interface{}) (bool, []interface{}) {
	if searches == nil {
		return true, comments
	}
	var matches []interface{}
	for _, search := range searches {
		if search.Reductive {
			var newMatches []interface{}
			for _, match := range matches {
				if com, ok := match.(Comment); ok {
					anyContent, allContent := searchContent(search.Content, com.Message, search.CaseSensitive)
					anyUser, allUser := searchContent(search.Users, com.AuthorName, search.CaseSensitive)
					if (search.MatchAny && (anyContent || anyUser)) || allContent && allUser {
						newMatches = append(newMatches, com)
					} else {
						comments = append(comments, com)
					}
					continue
				}
				if rep, ok := match.(Reply); ok {
					anyContent, allContent := searchContent(search.Content, rep.Message, search.CaseSensitive)
					anyUser, allUser := searchContent(search.Users, rep.AuthorName, search.CaseSensitive)
					if (search.MatchAny && (anyContent || anyUser)) || allContent && allUser {
						newMatches = append(newMatches, rep)
					} else {
						comments = append(comments, rep)
					}
					continue
				}
				return false, nil
			}
			matches = newMatches
		} else {
			var commentsLeft []interface{}
			for _, comment := range comments {
				if com, ok := comment.(Comment); ok {
					anyContent, allContent := searchContent(search.Content, com.Message, search.CaseSensitive)
					anyUser, allUser := searchContent(search.Users, com.AuthorName, search.CaseSensitive)
					if (search.MatchAny && (anyContent || anyUser)) || allContent && allUser {
						matches = append(matches, com)
					} else {
						commentsLeft = append(commentsLeft, com)
					}
					continue
				}
				if rep, ok := comment.(Reply); ok {
					anyContent, allContent := searchContent(search.Content, rep.Message, search.CaseSensitive)
					anyUser, allUser := searchContent(search.Users, rep.AuthorName, search.CaseSensitive)
					if (search.MatchAny && (anyContent || anyUser)) || allContent && allUser {
						matches = append(matches, rep)
					} else {
						commentsLeft = append(commentsLeft, rep)
					}
					continue
				}
				return false, nil
			}
			comments = commentsLeft
		}
	}
	return true, matches
}

func CommentsHandler(input Inputs) http.Handler {
	comments := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			key := r.URL.Query().Get("key")
			if key == "" {
				sendStatusCode(w, http.StatusBadRequest, "keyMissing")
				return
			}
			id := r.URL.Query().Get("id")
			if id == "" {
				sendStatusCode(w, http.StatusBadRequest, "videoIdMissing")
				return
			}
			var searches []Search
			if r.Body != nil {
				r.Body = http.MaxBytesReader(w, r.Body, 1048576) // Read max 1 MB
				searchErr := json.NewDecoder(r.Body).Decode(&searches)
				if searchErr != nil && searchErr.Error() == "http: request body too large" {
					sendStatusCode(w, http.StatusRequestEntityTooLarge, "searchBodyTooLarge")
					return
				} else if searchErr != nil && searchErr != io.EOF {
					print(searchErr.Error())
					sendStatusCode(w, http.StatusBadRequest, "searchBodyInvalid")
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
				ok:= func() bool {  // Internal function for deferring the closing of response bodies inside loop.
					resp, err := http.Get(fmt.Sprintf("%s&videoId=%s&key=%s&pageToken=%s",
						input.CommentsRoot, id, key, pageToken))
					if err != nil {
						sendStatusCode(w, http.StatusInternalServerError, "failedToQueryYouTubeAPI")
						return false
					}
					defer resp.Body.Close()
					youtubeStatus = ErrorParser(resp.Body, &commentsInbound)
					if youtubeStatus.StatusCode != http.StatusOK {
						sendStatusCode(w, youtubeStatus.StatusCode, youtubeStatus.StatusMessage)
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
			for _, comId := range needReplies {
				pageToken = ""
				for hasNextPage := true; hasNextPage; hasNextPage = pageToken != "" {
					var youtubeStatus StatusCodeOutbound
					var repliesInbound RepliesInbound
					ok:= func() bool {  // Internal function for deferring the closing of response bodies inside loop.
						resp, err := http.Get(fmt.Sprintf("%s&parentId=%s&key=%s&pageToken=%s",
							input.RepliesRoot, comId, key, pageToken))
						if err != nil {
							sendStatusCode(w, http.StatusInternalServerError, "failedToQueryYouTubeAPI")
							return false
						}
						defer resp.Body.Close()
						youtubeStatus = ErrorParser(resp.Body, &repliesInbound)
						if youtubeStatus.StatusCode != http.StatusOK {
							sendStatusCode(w, youtubeStatus.StatusCode, youtubeStatus.StatusMessage)
							return false
						}
						RepliesParser(repliesInbound, &comments)
						pageToken = repliesInbound.NextPageToken
						return true
					}
					if !ok() {
						return
					}
				}
			}
			ok, filteredComments := CommentSearch(searches, comments)
			if !ok {
				sendStatusCode(w, http.StatusInternalServerError, "failedFilteringComments")
				return
			}
			SortComments(&comments)
			commentsOutbound.Comments = filteredComments
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
