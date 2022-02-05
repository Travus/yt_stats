package yt_stats

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// ChatHandler is the handler for the chat endpoint. /ytstats/v1/chat/
// Lists messages of ongoing live streams. Only works on currently active streams.
func ChatHandler(input Inputs) http.Handler {
	stats := func(w http.ResponseWriter, r *http.Request) {
		quota := 0
		switch r.Method {
		case http.MethodGet:

			// Check user input and fail if input is incorrect or missing.
			var youtubeStatus StatusCodeOutbound
			var chatInbound ChatInbound
			key := getKey(r)
			if key == "" {
				sendStatusCode(w, quota, http.StatusBadRequest, "keyMissing")
				return
			}
			id := r.URL.Query().Get("id")
			if id == "" {
				sendStatusCode(w, quota, http.StatusBadRequest, "chatIdMissing")
				return
			}
			if len(strings.Split(id, ",")) > 1 {
				sendStatusCode(w, quota, http.StatusBadRequest, "tooManyItems")
				return
			}

			// Query youtube and check response for errors.
			resp, err := http.Get(fmt.Sprintf("%s&liveChatId=%s&key=%s", input.ChatRoot, url.QueryEscape(id), key))
			if err != nil {
				sendStatusCode(w, quota, http.StatusInternalServerError, "failedToQueryYouTubeAPI")
				return
			}
			defer resp.Body.Close()
			quota += 5
			youtubeStatus = ErrorParser(resp.Body, &chatInbound)
			if youtubeStatus.StatusCode != http.StatusOK {
				if youtubeStatus.StatusMessage == "keyInvalid" { // Quota cannot be deducted from invalid keys.
					quota -= 5
				}
				sendStatusCode(w, quota, youtubeStatus.StatusCode, youtubeStatus.StatusMessage)
				return
			}

			// Process and provide response.
			chatOutbound := ChatParser(chatInbound, id)
			chatOutbound.QuotaUsage = quota
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(chatOutbound)
			if err != nil {
				log.Println("Failed to respond to chat endpoint.")
			}
			return
		default:
			unsupportedRequestType(w)
			return
		}
	}
	return http.HandlerFunc(stats)
}
