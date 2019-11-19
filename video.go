package yt_stats

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func VideoHandler(input Inputs) http.Handler {
	video := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			var youtubeStatus StatusCodeOutbound
			videoInbound := make([]VideoInbound, 1)
			key := r.URL.Query().Get("key")
			if key == "" {
				sendStatusCode(w, http.StatusBadRequest, "keyMissing")
				return
			}
			ids := r.URL.Query().Get("id")
			if ids == "" {
				sendStatusCode(w, http.StatusBadRequest, "channelIdMissing")
				return
			}
			if len(strings.Split(ids, ",")) > 50 {
				sendStatusCode(w, http.StatusBadRequest, "tooManyItems")
				return
			}
			statsFlag := strings.ToLower(r.URL.Query().Get("stats"))
			if statsFlag != "" && statsFlag != "true" && statsFlag != "false" {
				sendStatusCode(w, http.StatusBadRequest, "flagInvalid")
				return
			}
			resp, err := http.Get(fmt.Sprintf("%s?part=snippet,contentDetails,statistics&id=%s&key=%s&" +
				"maxResults=50", input.VideosRoot, url.QueryEscape(ids), key))
			if err != nil {
				sendStatusCode(w, http.StatusInternalServerError, "failedToQueryYouTubeAPI")
				return
			}
			defer resp.Body.Close()
			youtubeStatus = ErrorParser(resp.Body, &videoInbound[0])
			if youtubeStatus.StatusCode != http.StatusOK {
				sendStatusCode(w, youtubeStatus.StatusCode, youtubeStatus.StatusMessage)
				return
			}
			var tempPlaylistObject Playlist
			var videoOutbound VideoOutbound
			err = VideoParser(videoInbound, &tempPlaylistObject, statsFlag == "true", true)
			if err != nil {
				sendStatusCode(w, http.StatusInternalServerError, "failedParsingYouTubeResponse")
			}
			videoOutbound.VideoStats = tempPlaylistObject.VideoStats
			videoOutbound.Videos = tempPlaylistObject.Videos
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
