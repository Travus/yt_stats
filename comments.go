package yt_stats

import (
	"net/http"
)

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
