package yt_stats

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func StatusHandler(input Inputs) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			var youtubeStatus StatusCodeOutbound
			key := r.URL.Query().Get("key")
			if key == "" {
				youtubeStatus = StatusCodeOutbound{
					StatusCode:    http.StatusBadRequest,
					StatusMessage: "keyMissing",
				}
				w.WriteHeader(http.StatusBadRequest)
				err := json.NewEncoder(w).Encode(youtubeStatus)
				if err != nil {
					log.Println("Failed to respond to status endpoint.")
				}
				return
			}
			uptime := time.Since(input.StartTime).Round(time.Second).Seconds()
			resp, err := http.Get(fmt.Sprintf("%s?part=id&id=UCBR8-60-B28hp2BmDPdntcQ&key=%s",
				input.ChannelsRoot, key))
			if err != nil {
				youtubeStatus = StatusCodeOutbound{
					StatusCode:    http.StatusInternalServerError,
					StatusMessage: "yt_stats API failed to query YouTube"}
			} else {
				youtubeStatus = ErrorParser(resp.Body, nil)
			}
			if resp != nil {
				defer resp.Body.Close()
			}
			response := StatusOutbound{Version: "v1", Uptime: uptime, YoutubeStatus: youtubeStatus}
			err = json.NewEncoder(w).Encode(response)
			if err != nil {
				log.Println("Failed to respond to status endpoint.")
			}
			break
		default:
			http.Error(w, "Request type not supported.", http.StatusNotImplemented)
		}
	})
}
