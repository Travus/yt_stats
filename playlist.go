package yt_stats

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// PlaylistHandler is the handler for the playlist endpoint. /ytstats/v1/playlist/
// Provides information on one or more playlists, and optionally a list of video and statistics on them.
func PlaylistHandler(input Inputs) http.Handler {
	playlist := func(w http.ResponseWriter, r *http.Request) {
		quota := 0
		switch r.Method {
		case http.MethodGet:

			// Check user input and fail if input is incorrect or missing.
			var youtubeStatus StatusCodeOutbound
			var playlistInbound PlaylistInbound
			key := r.URL.Query().Get("key")
			if key == "" {
				sendStatusCode(w, quota, http.StatusBadRequest, "keyMissing")
				return
			}
			ids := r.URL.Query().Get("id")
			if ids == "" {
				sendStatusCode(w, quota, http.StatusBadRequest, "playlistIdMissing")
				return
			}
			if len(strings.Split(ids, ",")) > 50 {
				sendStatusCode(w, quota, http.StatusBadRequest, "tooManyItems")
				return
			}
			videosFlag := strings.ToLower(r.URL.Query().Get("videos"))
			if videosFlag != "" && videosFlag != "true" && videosFlag != "false" {
				sendStatusCode(w, quota, http.StatusBadRequest, "flagInvalid")
				return
			}
			statsFlag := strings.ToLower(r.URL.Query().Get("stats"))
			if statsFlag != "" && statsFlag != "true" && statsFlag != "false" {
				sendStatusCode(w, quota, http.StatusBadRequest, "flagInvalid")
				return
			}

			// Query youtube playlist endpoint and check response for errors.
			resp, err := http.Get(fmt.Sprintf("%s&id=%s&key=%s", input.PlaylistsRoot, url.QueryEscape(ids), key))
			if err != nil {
				sendStatusCode(w, quota, http.StatusInternalServerError, "failedToQueryYouTubeAPI")
				return
			}
			defer resp.Body.Close()
			quota += 5
			youtubeStatus = ErrorParser(resp.Body, &playlistInbound)
			if youtubeStatus.StatusCode != http.StatusOK {
				if youtubeStatus.StatusMessage == "keyInvalid" {  // Quota cannot be deducted from invalid keys.
					quota -= 5
				}
				sendStatusCode(w, quota, youtubeStatus.StatusCode, youtubeStatus.StatusMessage)
				return
			}

			// Parse response and provide response if no video or statistics info was wanted.
			plOutbound := PlaylistTopLevelParser(playlistInbound)
			if videosFlag == "false" && statsFlag == "false" {
				plOutbound.QuotaUsage = quota
				w.Header().Set("Content-Type", "application/json")
				err = json.NewEncoder(w).Encode(plOutbound)
				if err != nil {
					log.Println("Failed to respond to playlist endpoint.")
				}
				return
			}

			// For all playlists query playlist items endpoint and handle pagination.
			for i := range plOutbound.Playlists {
				pageToken := ""
				var playlistItemsInbound []PlaylistItemsInbound
				for hasNextPage := true; hasNextPage; hasNextPage = pageToken != "" {
					var playlistItemPageInbound PlaylistItemsInbound
					ok := func() bool { // Internal function for deferring the closing of response bodies inside loop.
						resp, err = http.Get(fmt.Sprintf("%s&playlistId=%s&key=%s&pageToken=%s",
							input.PlaylistItemsRoot, plOutbound.Playlists[i].Id, key, pageToken))
						if err != nil {
							sendStatusCode(w, quota, http.StatusInternalServerError, "failedToQueryYouTubeAPI")
							return false
						}
						defer resp.Body.Close()
						quota += 3
						youtubeStatus = ErrorParser(resp.Body, &playlistItemPageInbound)
						if youtubeStatus.StatusCode != http.StatusOK {
							sendStatusCode(w, quota, youtubeStatus.StatusCode, youtubeStatus.StatusMessage)
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

				// Parse video IDs from playlist items endpoint response, query youtube videos endpoint and parse.
				videoIds := PlaylistItemsParser(playlistItemsInbound)
				var videoInbound []VideoInbound
				for _, page := range videoIds {
					var videoInboundPage VideoInbound
					ok := func() bool {
						videoPageIds := url.QueryEscape(strings.Join(page, ","))
						resp, err = http.Get(fmt.Sprintf("%s&id=%s&key=%s", input.VideosRoot, videoPageIds, key))
						if err != nil {
							sendStatusCode(w, quota, http.StatusInternalServerError, "failedToQueryYouTubeAPI")
							return false
						}
						defer resp.Body.Close()
						quota += 7
						youtubeStatus = ErrorParser(resp.Body, &videoInboundPage)
						if youtubeStatus.StatusCode != http.StatusOK {
							sendStatusCode(w, quota, youtubeStatus.StatusCode, youtubeStatus.StatusMessage)
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
					sendStatusCode(w, quota, http.StatusInternalServerError, "failedParsingYouTubeResponse")
				}
			}

			// Provide response.
			plOutbound.QuotaUsage = quota
			w.Header().Set("Content-Type", "application/json")
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
