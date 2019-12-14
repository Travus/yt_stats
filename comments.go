package yt_stats

import (
	"fmt"
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
	var matches []interface{}
	for _, search := range searches {
		fmt.Printf("%+v", search)
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
		fmt.Printf("matches: %d, comments: %d\n", len(matches), len(comments))
	}
	return true, matches
}

func CommentsHandler(input Inputs) http.Handler {
	comments := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			return
		default:
			unsupportedRequestType(w)
			return
		}
	}
	return http.HandlerFunc(comments)
}
