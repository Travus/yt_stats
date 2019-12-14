package yt_stats

import (
	"net/http"
	"sort"
)

// Sorts the comments slice based on comment and reply publishing time.
func SortComments(comments *[]interface{}) {
	sortFunc := func(i int, j int) bool {
		if com1, cok1 := (*comments)[i].(Comment); cok1 {
			if com2, cok2 := (*comments)[j].(Comment); cok2 {
				return com1.PublishedAt < com2.PublishedAt
			} else if rep2, rok2 := (*comments)[j].(Reply); rok2 {
				return com1.PublishedAt < rep2.PublishedAt
			}
		} else if rep1, rok1 := (*comments)[i].(Reply); rok1 {
			if com2, cok2 := (*comments)[j].(Comment); cok2 {
				return rep1.PublishedAt < com2.PublishedAt
			} else if rep2, rok2 := (*comments)[j].(Reply); rok2 {
				return rep1.PublishedAt < rep2.PublishedAt
			}
		}
		return false
	}
	sort.Slice(*comments, sortFunc)
}

func CommentSearch(searches []Search, matches *[]Comment, others *[]Comment) error {
	return nil
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
