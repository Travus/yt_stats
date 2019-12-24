package yt_stats

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func StatusHandler(input Inputs) http.Handler {
	stats := func(w http.ResponseWriter, r *http.Request) {
		quota := 0
		switch r.Method {
		case http.MethodGet:
			var youtubeStatus StatusCodeOutbound
			key := r.URL.Query().Get("key")
			if key == "" {
				sendStatusCode(w, quota, http.StatusBadRequest, "keyMissing")
				return
			}
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
