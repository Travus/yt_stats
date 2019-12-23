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
			key := r.URL.Query().Get("key")
			if key == "" {
				sendStatusCode(w, http.StatusBadRequest, "keyMissing")
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
			resp, err := http.Get(fmt.Sprintf("%s&id=%s&key=%s",
				input.PlaylistsRoot, url.QueryEscape(ids), key))
			if err != nil {
				sendStatusCode(w, http.StatusInternalServerError, "failedToQueryYouTubeAPI")
				return
			}
			defer resp.Body.Close()
			youtubeStatus = ErrorParser(resp.Body, &playlistInbound)
			if youtubeStatus.StatusCode != http.StatusOK {
				sendStatusCode(w, youtubeStatus.StatusCode, youtubeStatus.StatusMessage)
				return
			}
			plOutbound := PlaylistTopLevelParser(playlistInbound)
			if videosFlag == "false" && statsFlag == "false" {
				err = json.NewEncoder(w).Encode(plOutbound)
				if err != nil {
					log.Println("Failed to respond to playlist endpoint.")
				}
				return
			}
			for i := range plOutbound.Playlists {
				pageToken := ""
				var playlistItemsInbound []PlaylistItemsInbound
				for hasNextPage := true; hasNextPage; hasNextPage = pageToken != "" {
					var playlistItemPageInbound PlaylistItemsInbound
					ok := func() bool {  // Internal function for deferring the closing of response bodies inside loop.
						resp, err = http.Get(fmt.Sprintf("%s&playlistId=%s&key=%s&pageToken=%s",
							input.PlaylistItemsRootRoot, plOutbound.Playlists[i].Id, key, pageToken))
						if err != nil {
							sendStatusCode(w, http.StatusInternalServerError, "failedToQueryYouTubeAPI")
							return false
						}
						defer resp.Body.Close()
						youtubeStatus = ErrorParser(resp.Body, &playlistItemPageInbound)
						if youtubeStatus.StatusCode != http.StatusOK {
							sendStatusCode(w, youtubeStatus.StatusCode, youtubeStatus.StatusMessage)
							return false
						}
						pageToken = playlistItemPageInbound.NextPageToken
						playlistItemsInbound = append(playlistItemsInbound, playlistItemPageInbound)
						return true
					}
					if !ok() {
						return
					}
				}
				videoIds := PlaylistItemsParser(playlistItemsInbound)
				var videoInbound []VideoInbound
				for _, page := range videoIds {
					var videoInboundPage VideoInbound
					ok := func() bool {
						videoPageIds := url.QueryEscape(strings.Join(page, ","))
						resp, err = http.Get(fmt.Sprintf("%s&id=%s&key=%s", input.VideosRoot, videoPageIds, key))
						if err != nil {
							sendStatusCode(w, http.StatusInternalServerError, "failedToQueryYouTubeAPI")
							return false
						}
						defer resp.Body.Close()
						youtubeStatus = ErrorParser(resp.Body, &videoInboundPage)
						if youtubeStatus.StatusCode != http.StatusOK {
							sendStatusCode(w, youtubeStatus.StatusCode, youtubeStatus.StatusMessage)
							return false
						}
						videoInbound = append(videoInbound, videoInboundPage)
						return true
					}()
					if !ok {
						return
					}
				}
				err = VideoParser(videoInbound, &plOutbound.Playlists[i], statsFlag != "false", videosFlag != "false")
				if err != nil {
					sendStatusCode(w, http.StatusInternalServerError, "failedParsingYouTubeResponse")
				}
			}
			err = json.NewEncoder(w).Encode(plOutbound)
			if err != nil {
				log.Println("Failed to respond to playlist endpoint.")
			}
			return
		default:
			unsupportedRequestType(w)
			return
		}
	}
	return http.HandlerFunc(playlist)
}
