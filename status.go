package yt_stats

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		key := r.URL.Query().Get("token")
		var youtubeStatus StatusCodeOutbound
		uptime := time.Since(StartTime).Round(time.Second).Seconds()
		resp, err := http.Get(fmt.Sprintf("%s?part=id&id=UCBR8-60-B28hp2BmDPdntcQ&key=%s", ChannelsRoot, key))
		if err != nil {
			youtubeStatus = StatusCodeOutbound{StatusCode: http.StatusInternalServerError,
				                               StatusMessage: "yt_stats API failed to query YouTube"}
		} else {
			youtubeStatus = ErrorParser(resp.Body, nil)
		}
		defer resp.Body.Close()
		response := StatusOutbound{Version: "v1", Uptime: uptime, YoutubeStatus: youtubeStatus}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Println("Failed to respond to status endpoint.")
		}
		break
	default:
		http.Error(w, "Request type not supported.", http.StatusNotImplemented)
	}
}
