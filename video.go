package yt_stats

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// VideoHandler is the handler for the video endpoint. /ytstats/v1/video/
// Provides info of up to 50 videos, and statistics of them.
func VideoHandler(input Inputs) http.Handler {
	video := func(w http.ResponseWriter, r *http.Request) {
		quota := 0
		switch r.Method {
		case http.MethodGet:

			// Check user input and fail if input is incorrect or missing.
			var youtubeStatus StatusCodeOutbound
			videoInbound := make([]VideoInbound, 1)
			key := r.URL.Query().Get("key")
			if key == "" {
				sendStatusCode(w, quota, http.StatusBadRequest, "keyMissing")
				return
			}
			ids := r.URL.Query().Get("id")
			if ids == "" {
				sendStatusCode(w, quota, http.StatusBadRequest, "videoIdMissing")
				return
			}
			if len(strings.Split(ids, ",")) > 50 {
				sendStatusCode(w, quota, http.StatusBadRequest, "tooManyItems")
				return
			}
			statsFlag := strings.ToLower(r.URL.Query().Get("stats"))
			if statsFlag != "" && statsFlag != "true" && statsFlag != "false" {
				sendStatusCode(w, quota, http.StatusBadRequest, "flagInvalid")
				return
			}

			// Query youtube videos endpoint and handle errors.
			resp, err := http.Get(fmt.Sprintf("%s&id=%s&key=%s", input.VideosRoot, url.QueryEscape(ids), key))
			if err != nil {
				sendStatusCode(w, quota, http.StatusInternalServerError, "failedToQueryYouTubeAPI")
				return
			}
			defer resp.Body.Close()
			quota++
			youtubeStatus = ErrorParser(resp.Body, &videoInbound[0])
			if youtubeStatus.StatusCode != http.StatusOK {
				if youtubeStatus.StatusMessage == "keyInvalid" { // Quota cannot be deducted from invalid keys.
					quota--
				}
				sendStatusCode(w, quota, youtubeStatus.StatusCode, youtubeStatus.StatusMessage)
				return
			}

			// Parse videos endpoint response and provide response.
			var tempPlaylistObject Playlist
			var videoOutbound VideoOutbound
			err = VideoParser(videoInbound, &tempPlaylistObject, statsFlag == "true", true)
			if err != nil {
				sendStatusCode(w, quota, http.StatusInternalServerError, "failedParsingYouTubeResponse")
			}
			videoOutbound.VideoStats = tempPlaylistObject.VideoStats
			videoOutbound.Videos = tempPlaylistObject.Videos
			videoOutbound.QuotaUsage = quota
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(videoOutbound)
			if err != nil {
				log.Println("Failed to respond to playlist endpoint.")
			}
			return
		default:
			unsupportedRequestType(w)
			return
		}
	}
	return http.HandlerFunc(video)
}
