package yt_stats

import (
	"net/http"
)

func PlaylistHandler(input Inputs) http.Handler {
	playlist := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			return
		default:
			unsupportedRequestType(w)
			return
		}
	}
	return http.HandlerFunc(playlist)
}
