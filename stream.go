package yt_stats

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// StreamHandler is the handler for the stream endpoint. /ytstats/v1/stream/
// Provides status and information on live streams such as start time, and if currently online concurrent viewers.
func StreamHandler(input Inputs) http.Handler {
	stats := func(w http.ResponseWriter, r *http.Request) {
		quota := 0
		switch r.Method {
		case http.MethodGet:

			// Check user input and fail if input is incorrect or missing.
			var youtubeStatus StatusCodeOutbound
			var streamInbound StreamInbound
			key := r.URL.Query().Get("key")
			if key == "" {
				sendStatusCode(w, quota, http.StatusBadRequest, "keyMissing")
				return
			}
			ids := r.URL.Query().Get("id")
			if ids == "" {
				sendStatusCode(w, quota, http.StatusBadRequest, "streamIdMissing")
				return
			}
			if len(strings.Split(ids, ",")) > 50 {
				sendStatusCode(w, quota, http.StatusBadRequest, "tooManyItems")
				return
			}

			// Query youtube and check response for errors.
			resp, err := http.Get(fmt.Sprintf("%s&id=%s&key=%s", input.StreamRoot, url.QueryEscape(ids), key))
			if err != nil {
				sendStatusCode(w, quota, http.StatusInternalServerError, "failedToQueryYouTubeAPI")
				return
			}
			defer resp.Body.Close()
			quota += 3
			youtubeStatus = ErrorParser(resp.Body, &streamInbound)
			if youtubeStatus.StatusCode != http.StatusOK {
				if youtubeStatus.StatusMessage == "keyInvalid" { // Quota cannot be deducted from invalid keys.
					quota -= 3
				}
				sendStatusCode(w, quota, youtubeStatus.StatusCode, youtubeStatus.StatusMessage)
				return
			}

			// Process and provide response.
			streamOutbound := StreamParser(streamInbound)
			streamOutbound.QuotaUsage = quota
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(streamOutbound)
			if err != nil {
				log.Println("Failed to respond to stream endpoint.")
			}
			return
		default:
			unsupportedRequestType(w)
			return
		}
	}
	return http.HandlerFunc(stats)
}
