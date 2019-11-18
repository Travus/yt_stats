package yt_stats

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func ChannelHandler(input Inputs) http.Handler {
	channel := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			var youtubeStatus StatusCodeOutbound
			var channelInbound ChannelInbound
			key := getKey(w, r)
			if key == "" {
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
			resp, err := http.Get(fmt.Sprintf("%s?part=id,snippet,contentDetails,statistics&id=%s&key=%s&maxResults=50",
				input.ChannelsRoot, url.QueryEscape(ids), key))
			if err != nil {
				sendStatusCode(w, http.StatusInternalServerError, "failedToQueryYouTubeAPI")
				return
			}
			defer resp.Body.Close()
			youtubeStatus = ErrorParser(resp.Body, &channelInbound)
			if youtubeStatus.StatusCode != http.StatusOK {
				sendStatusCode(w, youtubeStatus.StatusCode, youtubeStatus.StatusMessage)
				return
			}
			channelOutbound := ChannelParser(channelInbound)
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
