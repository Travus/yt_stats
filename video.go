package yt_stats

import "net/http"

func VideoHandler(input Inputs) http.Handler {
	video := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			return
		default:
			unsupportedRequestType(w)
			return
		}
	}
	return http.HandlerFunc(video)
}
