package yt_stats

import (
	"net/http"
)

// Handler for the stream endpoint. /ytstats/v1/stream/
// Provides status and information on live streams such as start time, and if currently online concurrent viewers.
func StreamHandler(input Inputs) http.Handler {
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