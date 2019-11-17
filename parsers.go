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
	var outbound []string
	for _, inboundPlItems := range inbound {
		for _, plItem := range inboundPlItems.Items {
			outbound = append(outbound, plItem.Snippet.ResourceId.VideoId)
		}
	}
	return outbound
}

func VideoParser(inbound []VideoInbound, playlistObject *Playlist, stats bool, videos bool) error {
	var vStats VideoStats
	totalViews, totalLikes, totalDislikes, totalComments := 0, 0, 0, 0
	for _, videoInbound := range inbound {
		for _, video := range videoInbound.Items {
			dur, err := durationConverter(video.ContentDetails.Duration)
			if err != nil {
				return err
			}
			views, err := strconv.Atoi(video.Statistics.ViewCount)
			if err != nil {
				return err
			}
			totalViews += views
			likes, err := strconv.Atoi(video.Statistics.LikeCount)
			if err != nil {
				return err
			}
			totalLikes += likes
			dislikes, err := strconv.Atoi(video.Statistics.DislikeCount)
			if err != nil {
				return err
			}
			totalDislikes += dislikes
			comments, err := strconv.Atoi(video.Statistics.CommentCount)
			if err != nil {
				return err
			}
			totalComments += comments
			if videos {
				vid := Video{
					Id:           video.Id,
					Title:        video.Snippet.Title,
					Description:  video.Snippet.Description,
					PublishedAt:  video.Snippet.PublishedAt,
					Thumbnail:    video.Snippet.Thumbnails.Medium.Url,
					ChannelId:    video.Snippet.ChannelId,
					Duration:     dur,
					ViewCount:    views,
					LikeCount:    likes,
					DislikeCount: dislikes,
					CommentCount: comments,
				}
				playlistObject.Videos = append(playlistObject.Videos, vid)
			}
			if stats {
				vStats.AvailableVideos++
				vStats.TotalLength += dur
				vStats.TotalViews += views
				if dur > vStats.LongestVideoDuration || vStats.LongestVideo == "" {
					vStats.LongestVideo = video.Id
					vStats.LongestVideoDuration = dur
				}
				if dur < vStats.ShortestVideoDuration || vStats.ShortestVideo == "" {
					vStats.ShortestVideo = video.Id
					vStats.ShortestVideoDuration = dur
				}
				if views > vStats.MostViews || vStats.MostViewedVideo == "" {
					vStats.MostViewedVideo = video.Id
					vStats.MostViews = views
				}
				if views < vStats.LeastViews || vStats.LeastViewedVideo == "" {
					vStats.LeastViewedVideo = video.Id
					vStats.LeastViews = views
				}
				if likes > vStats.MostLikes || vStats.MostLikedVideo == "" {
					vStats.MostLikedVideo = video.Id
					vStats.MostLikes = likes
				}
				if likes < vStats.LeastLikes || vStats.LeastLikedVideo == "" {
					vStats.LeastLikedVideo = video.Id
					vStats.LeastLikes = likes
				}
				if dislikes > vStats.MostDislikes || vStats.MostDislikedVideo == "" {
					vStats.MostDislikedVideo = video.Id
					vStats.MostDislikes = dislikes
				}
				if dislikes < vStats.LeastDislikes || vStats.LeastDislikedVideo == "" {
					vStats.LeastDislikedVideo = video.Id
					vStats.LeastDislikes = dislikes
				}
				if comments > vStats.MostComments || vStats.MostCommentedVideo == "" {
					vStats.MostCommentedVideo = video.Id
					vStats.MostComments = comments
				}
				if comments < vStats.LeastComments || vStats.LeastCommentedVideo == "" {
					vStats.LeastCommentedVideo = video.Id
					vStats.LeastComments = comments
				}
			}
		}
	}
	if stats {
		vStats.AverageVideoDuration = vStats.TotalLength / vStats.AvailableVideos
		vStats.AverageViews = vStats.TotalViews / vStats.AvailableVideos
		vStats.AverageLikes = totalLikes / vStats.AvailableVideos
		vStats.AverageDislikes = totalDislikes / vStats.AvailableVideos
		vStats.AverageComments = totalComments / vStats.AvailableVideos
		playlistObject.VideoStats = vStats
	}
	return nil
}

func FullPlaylistParser(inbound PlaylistInbound, stats bool, videos bool) PlaylistOutbound {
	return PlaylistOutbound{}
}
