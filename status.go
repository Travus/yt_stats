package yt_stats

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// StatusHandler is the handler for the status endpoint. /ytstats/v1/status/
// Provides version, uptime, and the status of the youtube API.
func StatusHandler(input Inputs) http.Handler {
	stats := func(w http.ResponseWriter, r *http.Request) {
		quota := 0
		switch r.Method {
		case http.MethodGet:

			// Check user input, this endpoint is allowed to progress even without a key.
			var youtubeStatus StatusCodeOutbound
			key := getKey(r)

			// Query youtube to check for youtube API status.
			uptime := time.Since(input.StartTime).Round(time.Second).Seconds()
			statusRoot := "https://www.googleapis.com/youtube/v3/channels?part=id"
			resp, err := http.Get(fmt.Sprintf("%s&id=UCBR8-60-B28hp2BmDPdntcQ&key=%s", statusRoot, key))
			if err != nil {
				sendStatusCode(w, quota, http.StatusInternalServerError, "failedToQueryYouTubeAPI")
				return
			}
			defer resp.Body.Close()
			quota++
			youtubeStatus = ErrorParser(resp.Body, nil)
			if youtubeStatus.StatusMessage == "keyInvalid" { // Quota cannot be deducted from invalid keys.
				quota--
			}

			// Create and provide response.
			youtubeStatus.QuotaUsage = quota
			response := StatusOutbound{
				QuotaUsage: quota,
				Version:    "v1",
				Uptime:     uptime,
				YoutubeStatus: struct {
					StatusCode    int    `json:"status_code"`
					StatusMessage string `json:"status_message"`
				}{StatusCode: youtubeStatus.StatusCode, StatusMessage: youtubeStatus.StatusMessage},
			}
			w.WriteHeader(response.YoutubeStatus.StatusCode)
			err = json.NewEncoder(w).Encode(response)
			if err != nil {
				log.Println("Failed to respond to status endpoint.")
			}
			return
		default:
			unsupportedRequestType(w)
			return
		}
	}
	return http.HandlerFunc(stats)
}
