package yt_stats

import "time"

type Inputs struct {
	StartTime    time.Time
	RepliesRoot  string
	CommentsRoot string
	ChannelsRoot string
}

type YoutubeErrorInbound struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Errors  []struct {
			Reason string `json:"reason"`
		} `json:"errors"`
	} `json:"error"`
}

type StatusCodeOutbound struct {
	StatusCode    int    `json:"status_code"`
	StatusMessage string `json:"status_message"`
}

type StatusOutbound struct {
	Version       string             `json:"version"`
	Uptime        float64            `json:"uptime"`
	YoutubeStatus StatusCodeOutbound `json:"youtube_status"`
}

type ChannelInbound struct {
	Items []struct {
		Id      string `json:"id"`
		Snippet struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Thumbnails  struct {
				Medium struct {
					Url string `json:"url"`
				} `json:"medium"`
			} `json:"thumbnails"`
			Country string `json:"country"`
		} `json:"snippet"`
		ContentDetails struct {
			RelatedPlaylists struct {
				Uploads string `json:"uploads"`
			} `json:"relatedPlaylists"`
		} `json:"contentDetails"`
		Statistics struct {
			ViewCount             string `json:"viewCount"`
			SubscriberCount       string `json:"subscriberCount"`
			HiddenSubscriberCount bool   `json:"hiddenSubscriberCount"`
			VideoCount            string `json:"videoCount"`
		} `json:"statistics"`
	} `json:"items"`
}

type Channel struct {
	Id                    string `json:"id"`
	Title                 string `json:"title"`
	Description           string `json:"description"`
	Thumbnail             string `json:"thumbnail"`
	Country               string `json:"country"`
	UploadsPlaylist       string `json:"uploads_playlist"`
	ViewCount             int    `json:"view_count"`
	HiddenSubscriberCount bool   `json:"hidden_subscriber_count"`
	SubscriberCount       int    `json:"subscriber_count"`
	VideoCount            int    `json:"video_count"`
}

type ChannelOutbound struct {
	Channels []Channel `json:"channels"`
}
