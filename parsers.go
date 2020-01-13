package yt_stats

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Parses a youtube response into a given struct, and returns a status code struct with OK or response.
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
	return StatusCodeOutbound{
		StatusCode:    200,
		StatusMessage: "OK",
	}
}

// Parses a ChannelInbound struct to a ChannelOutbound struct.
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

// Parses a PlaylistInbound struct to a PlaylistOutbound struct.
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

// Parses a slice of PlaylistItemInbound structs to a slice of string slices including all the video IDs.
func PlaylistItemsParser(inbound []PlaylistItemsInbound) [][]string {
	var outbound [][]string
	for _, inboundPlItems := range inbound {
		var page []string
		for _, plItem := range inboundPlItems.Items {
			page = append(page, plItem.Snippet.ResourceId.VideoId)
		}
		outbound = append(outbound, page)
	}
	return outbound
}

// Parses a slice of VideoInbound structs into a Playlist struct and returns any errors.
func VideoParser(inbound []VideoInbound, playlistObject *Playlist, stats bool, videos bool) error {
	var vStats VideoStats
	totalViews, totalLikes, totalDislikes, totalComments := 0, 0, 0, 0
	for _, videoInbound := range inbound {

		// Handle overall statistics, like total duration etc.
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

			// Handle video specific statistics.
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

			// Handle most/longest and least/shortest statistics.
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

	// Handle average statistics.
	if stats {
		vStats.AverageVideoDuration = vStats.TotalLength / vStats.AvailableVideos
		vStats.AverageViews = vStats.TotalViews / vStats.AvailableVideos
		vStats.AverageLikes = totalLikes / vStats.AvailableVideos
		vStats.AverageDislikes = totalDislikes / vStats.AvailableVideos
		vStats.AverageComments = totalComments / vStats.AvailableVideos
		playlistObject.VideoStats = &vStats
	} else {
		playlistObject.VideoStats = nil
	}
	return nil
}

// Parses a CommentsInbound struct into a slice of interfaces containing Comment and Reply structs.
func CommentsParser(inbound CommentsInbound, comments *[]interface{}, replies *[]string) {
	for _, item := range inbound.Items {
		com := Comment{
			Type:             "comment",
			Id:               item.Snippet.TopLevelComment.Id,
			AuthorName:       item.Snippet.TopLevelComment.Snippet.AuthorDisplayName,
			AuthorId:         item.Snippet.TopLevelComment.Snippet.AuthorChannelId.Value,
			AuthorChannelURL: item.Snippet.TopLevelComment.Snippet.AuthorChannelUrl,
			Message:          item.Snippet.TopLevelComment.Snippet.TextDisplay,
			Likes:            item.Snippet.TopLevelComment.Snippet.LikeCount,
			PublishedAt:      item.Snippet.TopLevelComment.Snippet.PublishedAt,
			ReplyCount:       item.Snippet.TotalReplyCount,
		}
		*comments = append(*comments, com)
		if len(item.Replies.Comments) == item.Snippet.TotalReplyCount {
			for _, repItem := range item.Replies.Comments {
				rep := Reply{
					Type:             "reply",
					Id:               repItem.Id,
					ParentId:         repItem.Id[:strings.IndexByte(repItem.Id, '.')],
					AuthorName:       repItem.Snippet.AuthorDisplayName,
					AuthorId:         repItem.Snippet.AuthorChannelId.Value,
					AuthorChannelURL: repItem.Snippet.AuthorChannelUrl,
					Message:          repItem.Snippet.TextDisplay,
					Likes:            repItem.Snippet.LikeCount,
					PublishedAt:      repItem.Snippet.PublishedAt,
				}
				*comments = append(*comments, rep)
			}
		} else {
			*replies = append(*replies, com.Id)
		}
	}
}

// Parses a RepliesInbound struct into a slice of interfaces containing Comment and Reply structs.
func RepliesParser(inbound RepliesInbound, comments *[]interface{}) {
	for _, item := range inbound.Items {
		rep := Reply{
			Type:             "reply",
			Id:               item.Id,
			ParentId:         item.Id[:strings.IndexByte(item.Id, '.')],
			AuthorName:       item.Snippet.AuthorDisplayName,
			AuthorId:         item.Snippet.AuthorChannelId.Value,
			AuthorChannelURL: item.Snippet.AuthorChannelUrl,
			Message:          item.Snippet.TextDisplay,
			Likes:            item.Snippet.LikeCount,
			PublishedAt:      item.Snippet.PublishedAt,
		}
		*comments = append(*comments, rep)
	}
}

// Parses a StreamInbound struct into a StreamOutbound struct.
func StreamParser(inbound StreamInbound) StreamOutbound {
	var outbound StreamOutbound
	outbound.Streams = make([]interface{}, len(inbound.Items))
	for i, video := range inbound.Items {
		if video.LiveStreamingDetails.ActualStartTime != "" && video.LiveStreamingDetails.ActualEndTime == "" {
			viewers, err := strconv.Atoi(video.LiveStreamingDetails.ConcurrentViewers)
			if err != nil {
				viewers = -1
				log.Printf("failed to convert concurrent viewers for %s", video.Id)
			}
			outbound.Streams[i] = LiveStream{
				Id:                 video.Id,
				Status:             "live",
				ScheduledStartTime: video.LiveStreamingDetails.ScheduledStartTime,
				StartTime:          video.LiveStreamingDetails.ActualStartTime,
				ConcurrentViewers:  viewers,
				ChatId:             video.LiveStreamingDetails.ActiveLiveChatId,
			}
		} else if video.LiveStreamingDetails.ActualEndTime != "" {
			outbound.Streams[i] = Stream{
				Id:                 video.Id,
				Status:             "ended",
				ScheduledStartTime: video.LiveStreamingDetails.ScheduledStartTime,
				StartTime:          video.LiveStreamingDetails.ActualStartTime,
				EndTime:            video.LiveStreamingDetails.ActualEndTime,
			}
		} else if video.LiveStreamingDetails.ScheduledStartTime != "" {
			outbound.Streams[i] = Stream{
				Id:                 video.Id,
				Status:             "scheduled",
				ScheduledStartTime: video.LiveStreamingDetails.ScheduledStartTime,
			}
		} else {
			outbound.Streams[i] = Stream{
				Id:                 video.Id,
				Status:             "video",
			}
		}
	}
	return outbound
}
