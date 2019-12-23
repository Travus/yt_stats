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
		switch r.Method {
		case http.MethodGet:
			var youtubeStatus StatusCodeOutbound
			key := r.URL.Query().Get("key")
			if key == "" {
				sendStatusCode(w, http.StatusBadRequest, "keyMissing")
				return
			}
			uptime := time.Since(input.StartTime).Round(time.Second).Seconds()
			statusRoot := "https://www.googleapis.com/youtube/v3/channels?part=id"
			resp, err := http.Get(fmt.Sprintf("%s&id=UCBR8-60-B28hp2BmDPdntcQ&key=%s", statusRoot, key))
			if err != nil {
				sendStatusCode(w, http.StatusInternalServerError, "failedToQueryYouTubeAPI")
				return
			}
			defer resp.Body.Close()
			youtubeStatus = ErrorParser(resp.Body, nil)
			response := StatusOutbound{
				Version:       "v1",
				Uptime:        uptime,
				YoutubeStatus: youtubeStatus,
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
