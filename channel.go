package yt_stats

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// ChannelHandler is the handler for the channel endpoint. /ytstats/v1/channel/
// Provides statistics for up to 50 channels.
func ChannelHandler(input Inputs) http.Handler {
	channel := func(w http.ResponseWriter, r *http.Request) {
		quota := 0
		switch r.Method {
		case http.MethodGet:

			// Check user input and fail if input is incorrect or missing.
			var youtubeStatus StatusCodeOutbound
			var channelInbound ChannelInbound
			key := r.URL.Query().Get("key")
			if key == "" {
				sendStatusCode(w, quota, http.StatusBadRequest, "keyMissing")
				return
			}
			ids := r.URL.Query().Get("id")
			if ids == "" {
				sendStatusCode(w, quota, http.StatusBadRequest, "channelIdMissing")
				return
			}
			if len(strings.Split(ids, ",")) > 50 {
				sendStatusCode(w, quota, http.StatusBadRequest, "tooManyItems")
				return
			}

			// Query youtube and check response for errors.
			resp, err := http.Get(fmt.Sprintf("%s&id=%s&key=%s", input.ChannelsRoot, url.QueryEscape(ids), key))
			if err != nil {
				sendStatusCode(w, quota, http.StatusInternalServerError, "failedToQueryYouTubeAPI")
				return
			}
			defer resp.Body.Close()
			quota += 7
			youtubeStatus = ErrorParser(resp.Body, &channelInbound)
			if youtubeStatus.StatusCode != http.StatusOK {
				if youtubeStatus.StatusMessage == "keyInvalid" {  // Quota cannot be deducted from invalid keys.
					quota -= 7
				}
				sendStatusCode(w, quota, youtubeStatus.StatusCode, youtubeStatus.StatusMessage)
				return
			}

			// Process and provide response.
			channelOutbound := ChannelParser(channelInbound)
			channelOutbound.QuotaUsage = quota
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(channelOutbound)
			if err != nil {
				log.Println("Failed to respond to channel endpoint.")
			}
			return
		default:
			unsupportedRequestType(w)
			return
		}
	}
	return http.HandlerFunc(channel)
}
