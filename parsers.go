package yt_stats

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
)

// This function parses a youtube response into a given struct, and returns a status code struct with OK or response.
func ErrorParser(r io.Reader, s interface{}) StatusCodeOutbound {
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)
	var errorCode YoutubeErrorInbound
	err := json.NewDecoder(tee).Decode(&errorCode)
	if err != nil {
		return StatusCodeOutbound{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: "failedToQueryYouTubeAPI",
		}
	}
	if errorCode.Error.Code != 0 {
		if errorCode.Error.Errors[0].Reason == "" {
			return StatusCodeOutbound{
				StatusCode:    errorCode.Error.Code,
				StatusMessage: errorCode.Error.Message,
			}
		}
		return StatusCodeOutbound{
			StatusCode:    errorCode.Error.Code,
			StatusMessage: errorCode.Error.Errors[0].Reason,
		}
	}
	if s != nil {
		err = json.NewDecoder(&buf).Decode(&s)
		if err != nil {
			return StatusCodeOutbound{
				StatusCode:    http.StatusInternalServerError,
				StatusMessage: "failedToQueryYouTubeAPI",
			}
		}
	}
	return StatusCodeOutbound{StatusCode: 200, StatusMessage: "OK"}
}

// This function parses a inbound channel struct to a outbound channel struct.
func ChannelParser(inbound ChannelInbound) ChannelOutbound {
	var outbound ChannelOutbound
	outbound.Channels = make([]Channel, len(inbound.Items))
	for i, rawChannel := range inbound.Items {
		outbound.Channels[i].Id = rawChannel.Id
		outbound.Channels[i].Title = rawChannel.Snippet.Title
		outbound.Channels[i].Description = rawChannel.Snippet.Description
		outbound.Channels[i].Thumbnail = rawChannel.Snippet.Thumbnails.Medium.Url
		outbound.Channels[i].Country = rawChannel.Snippet.Country
		outbound.Channels[i].UploadsPlaylist = rawChannel.ContentDetails.RelatedPlaylists.Uploads
		outbound.Channels[i].HiddenSubscriberCount = rawChannel.Statistics.HiddenSubscriberCount
		if !outbound.Channels[i].HiddenSubscriberCount {
			n, err := strconv.Atoi(rawChannel.Statistics.SubscriberCount)
			if err != nil {
				n = 0
				log.Printf("failed to convert subscriber count for %s", outbound.Channels[i].Id)
			}
			outbound.Channels[i].SubscriberCount = n
		}
		n, err := strconv.Atoi(rawChannel.Statistics.ViewCount)
		if err != nil {
			n = 0
			log.Printf("failed to convert view count for %s", outbound.Channels[i].Id)
		}
		outbound.Channels[i].ViewCount = n
		n, err = strconv.Atoi(rawChannel.Statistics.VideoCount)
		if err != nil {
			n = 0
			log.Printf("failed to convert video count for %s", outbound.Channels[i].Id)
		}
		outbound.Channels[i].VideoCount = n
	}
	return outbound
}

func PlaylistTopLevelParser(inbound PlaylistInbound) PlaylistOutbound {
	var outbound PlaylistOutbound
	outbound.Playlists = make([]Playlist, len(inbound.Items))
	for i, inboundPl := range inbound.Items {
		outbound.Playlists[i].Id = inboundPl.Id
		outbound.Playlists[i].Title = inboundPl.Snippet.Title
		outbound.Playlists[i].Description = inboundPl.Snippet.Description
		outbound.Playlists[i].PublishedAt = inboundPl.Snippet.PublishedAt
		outbound.Playlists[i].Thumbnail = inboundPl.Snippet.Thumbnails.Medium.Url
		outbound.Playlists[i].TotalVideos = inboundPl.ContentDetails.ItemCount
		outbound.Playlists[i].ChannelInfo.ChannelId = inboundPl.Snippet.ChannelId
		outbound.Playlists[i].ChannelInfo.ChannelTitle = inboundPl.Snippet.ChannelTitle
	}
	return outbound
}

func PlaylistItemsParser(inbound []PlaylistItemsInbound) []string {
	return []string{}
}

func VideoParser(inbound []VideoInbound, playlistObject *Playlist, stats bool, videos bool) {
	return
}

func FullPlaylistParser(inbound PlaylistInbound, stats bool, videos bool) PlaylistOutbound {
	return PlaylistOutbound{}
}
