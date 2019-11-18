package yt_stats

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func PlaylistHandler(input Inputs) http.Handler {
	playlist := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			var youtubeStatus StatusCodeOutbound
			var playlistInbound PlaylistInbound
			key := getKey(w, r)
			if key == "" {
				return
			}
			ids := r.URL.Query().Get("id")
			if ids == "" {
				sendStatusCode(w, http.StatusBadRequest, "playlistIdMissing")
				return
			}
			if len(strings.Split(ids, ",")) > 50 {
				sendStatusCode(w, http.StatusBadRequest, "tooManyItems")
				return
			}
			videosFlag := strings.ToLower(r.URL.Query().Get("videos"))
			if videosFlag != "" && videosFlag != "true" && videosFlag != "false" {
				sendStatusCode(w, http.StatusBadRequest, "flagInvalid")
				return
			}
			statsFlag := strings.ToLower(r.URL.Query().Get("stats"))
			if statsFlag != "" && statsFlag != "true" && statsFlag != "false" {
				sendStatusCode(w, http.StatusBadRequest, "flagInvalid")
				return
			}
			resp, err := http.Get(fmt.Sprintf("%s?part=snippet,contentDetails&id=%s&key=%s&maxResults=50",
				input.PlaylistsRoot, url.QueryEscape(ids), key))
			if err != nil {
				sendStatusCode(w, http.StatusInternalServerError, "failedToQueryYouTubeAPI")
				return
			}
			youtubeStatus = ErrorParser(resp.Body, &playlistInbound)
			defer resp.Body.Close()
			if youtubeStatus.StatusCode != http.StatusOK {
				sendStatusCode(w, youtubeStatus.StatusCode, youtubeStatus.StatusMessage)
				return
			}
			playlistOutbound := PlaylistTopLevelParser(playlistInbound)
			if videosFlag == "false" && statsFlag == "false" {
				err = json.NewEncoder(w).Encode(playlistOutbound)
				if err != nil {
					log.Println("Failed to respond to playlist endpoint.")
				}
				return
			}
			return
		default:
			unsupportedRequestType(w)
			return
		}
	}
	return http.HandlerFunc(playlist)
}
