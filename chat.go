package yt_stats

import (
	"net/http"
)

// Handler for the chat endpoint. /ytstats/v1/chat/
// Lists messages of ongoing live streams. Only works on currently active streams.
func ChatHandler(input Inputs) http.Handler {
	stats := func(w http.ResponseWriter, r *http.Request) {
		// quota := 0
		switch r.Method {
		case http.MethodGet:

		default:
			unsupportedRequestType(w)
			return
		}
	}
	return http.HandlerFunc(stats)
}

