package yt_stats

import (
	"time"
)

// Stores variables sent to the handlers, basically global variables.
type Inputs struct {
	StartTime         time.Time
	RepliesRoot       string
	CommentsRoot      string
	ChannelsRoot      string
	PlaylistsRoot     string
	PlaylistItemsRoot string
	VideosRoot        string
}

// Represents the JSON received from a YouTube error response.
type YoutubeErrorInbound struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Errors  []struct {
			Reason string `json:"reason"`
		} `json:"errors"`
	} `json:"error"`
}

// Represents the JSON sent on errors.
type StatusCodeOutbound struct {
	QuotaUsage    int    `json:"quota_usage"`
	StatusCode    int    `json:"status_code"`
	StatusMessage string `json:"status_message"`
}

// Represents the JSON sent by the Status endpoint.
type StatusOutbound struct {
	QuotaUsage    int     `json:"quota_usage"`
	Version       string  `json:"version"`
	Uptime        float64 `json:"uptime"`
	YoutubeStatus struct {
		StatusCode    int    `json:"status_code"`
		StatusMessage string `json:"status_message"`
	} `json:"youtube_status"`
}

// Represents the JSON received from the YouTube Channels endpoint.
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

// Represents the JSON for one channel. Part of ChannelOutbound struct.
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

// Represents the JSON sent by the Channel endpoint.
type ChannelOutbound struct {
	QuotaUsage int       `json:"quota_usage"`
	Channels   []Channel `json:"channels"`
}

// Represents the JSON received from the YouTube Playlists endpoint.
type PlaylistInbound struct {
	Items []struct {
		Id      string `json:"id"`
		Snippet struct {
			PublishedAt string `json:"publishedAt"`
			ChannelId   string `json:"channelId"`
			Title       string `json:"title"`
			Description string `json:"description"`
			Thumbnails  struct {
				Medium struct {
					Url string `json:"url"`
				} `json:"medium"`
			} `json:"thumbnails"`
			ChannelTitle string `json:"channelTitle"`
		} `json:"snippet"`
		ContentDetails struct {
			ItemCount int `json:"itemCount"`
		} `json:"contentDetails"`
	} `json:"items"`
}

// Represents the JSON received from the YouTube PlaylistsItems endpoint.
type PlaylistItemsInbound struct {
	NextPageToken string `json:"nextPageToken"`
	Items         []struct {
		Snippet struct {
			ResourceId struct {
				VideoId string `json:"videoId"`
			} `json:"resourceId"`
		} `json:"snippet"`
	} `json:"items"`
}

// Represents the JSON received from the YouTube Videos endpoint.
type VideoInbound struct {
	Items []struct {
		Id      string `json:"id"`
		Snippet struct {
			PublishedAt string `json:"publishedAt"`
			ChannelId   string `json:"channelId"`
			Title       string `json:"title"`
			Description string `json:"description"`
			Thumbnails  struct {
				Medium struct {
					Url string `json:"url"`
				} `json:"medium"`
			} `json:"thumbnails"`
		} `json:"snippet"`
		ContentDetails struct {
			Duration string `json:"duration"`
		} `json:"contentDetails"`
		Statistics struct {
			ViewCount    string `json:"viewCount"`
			LikeCount    string `json:"likeCount"`
			DislikeCount string `json:"dislikeCount"`
			CommentCount string `json:"commentCount"`
		} `json:"statistics"`
	} `json:"items"`
}

// Represents the JSON for stats over a range of videos. Part of Playlist struct.
type VideoStats struct {
	AvailableVideos       int    `json:"available_videos"`
	TotalLength           int    `json:"total_length"`
	TotalViews            int    `json:"total_views"`
	LongestVideo          string `json:"longest_video"`
	LongestVideoDuration  int    `json:"longest_video_duration"`
	ShortestVideo         string `json:"shortest_video"`
	ShortestVideoDuration int    `json:"shortest_video_duration"`
	AverageVideoDuration  int    `json:"average_video_duration"`
	MostViewedVideo       string `json:"most_viewed_video"`
	MostViews             int    `json:"most_views"`
	LeastViewedVideo      string `json:"least_viewed_video"`
	LeastViews            int    `json:"least_views"`
	AverageViews          int    `json:"average_views"`
	MostLikedVideo        string `json:"most_liked_video"`
	MostLikes             int    `json:"most_likes"`
	LeastLikedVideo       string `json:"least_liked_video"`
	LeastLikes            int    `json:"least_likes"`
	AverageLikes          int    `json:"average_likes"`
	MostDislikedVideo     string `json:"most_disliked_video"`
	MostDislikes          int    `json:"most_dislikes"`
	LeastDislikedVideo    string `json:"least_disliked_video"`
	LeastDislikes         int    `json:"least_dislikes"`
	AverageDislikes       int    `json:"average_dislikes"`
	MostCommentedVideo    string `json:"most_commented_video"`
	MostComments          int    `json:"most_comments"`
	LeastCommentedVideo   string `json:"least_commented_video"`
	LeastComments         int    `json:"least_comments"`
	AverageComments       int    `json:"average_comments"`
}

// Represents the JSON for one video. Part of Playlist struct.
type Video struct {
	Id           string `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	PublishedAt  string `json:"published_at"`
	Thumbnail    string `json:"thumbnail"`
	ChannelId    string `json:"channel_id"`
	Duration     int    `json:"duration"`
	ViewCount    int    `json:"view_count"`
	LikeCount    int    `json:"like_count"`
	DislikeCount int    `json:"dislike_count"`
	CommentCount int    `json:"comment_count"`
}

// Represents the JSON for one playlist. Part of PlaylistOutbound struct.
type Playlist struct {
	Id          string      `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	PublishedAt string      `json:"published_at"`
	Thumbnail   string      `json:"thumbnail"`
	TotalVideos int         `json:"total_videos"`
	VideoStats  *VideoStats `json:"video_stats,omitempty"`
	Videos      []Video     `json:"videos,omitempty"`
	ChannelInfo struct {
		ChannelId    string `json:"channel_id"`
		ChannelTitle string `json:"channel_title"`
	} `json:"channel_info"`
}

// Represents the JSON sent by the Playlist endpoint.
type PlaylistOutbound struct {
	QuotaUsage int        `json:"quota_usage"`
	Playlists  []Playlist `json:"playlists"`
}

// Represents the JSON sent by the Video endpoint.
type VideoOutbound struct {
	QuotaUsage int         `json:"quota_usage"`
	VideoStats *VideoStats `json:"video_stats,omitempty"`
	Videos     []Video     `json:"videos"`
}

// Represents the JSON received from the YouTube CommentThreads endpoint.
type CommentsInbound struct {
	NextPageToken string `json:"nextPageToken"`
	Items         []struct {
		Snippet struct {
			TopLevelComment struct {
				Id      string `json:"id"`
				Snippet struct {
					AuthorDisplayName string `json:"authorDisplayName"`
					AuthorChannelUrl  string `json:"authorChannelUrl"`
					AuthorChannelId   struct {
						Value string `json:"value"`
					} `json:"authorChannelId"`
					TextDisplay string `json:"textDisplay"`
					LikeCount   int    `json:"likeCount"`
					PublishedAt string `json:"publishedAt"`
				} `json:"snippet"`
			} `json:"topLevelComment"`
			TotalReplyCount int `json:"totalReplyCount"`
		} `json:"snippet"`
		Replies struct {
			Comments []struct {
				Id      string `json:"id"`
				Snippet struct {
					AuthorDisplayName string `json:"authorDisplayName"`
					AuthorChannelUrl  string `json:"authorChannelUrl"`
					AuthorChannelId   struct {
						Value string `json:"value"`
					} `json:"authorChannelId"`
					TextDisplay string `json:"textDisplay"`
					LikeCount   int    `json:"likeCount"`
					PublishedAt string `json:"publishedAt"`
				} `json:"snippet"`
			} `json:"comments"`
		} `json:"replies"`
	} `json:"items"`
}

// Represents the JSON received from the YouTube Comments endpoint.
type RepliesInbound struct {
	NextPageToken string `json:"nextPageToken"`
	Items         []struct {
		Id      string `json:"id"`
		Snippet struct {
			AuthorDisplayName string `json:"authorDisplayName"`
			AuthorChannelUrl  string `json:"authorChannelUrl"`
			AuthorChannelId   struct {
				Value string `json:"value"`
			} `json:"authorChannelId"`
			TextDisplay string `json:"textDisplay"`
			LikeCount   int    `json:"likeCount"`
			PublishedAt string `json:"publishedAt"`
		} `json:"snippet"`
	} `json:"items"`
}

// Represents the JSON for one comment. Part of CommentOutbound.
type Comment struct {
	Type             string `json:"type"`
	Id               string `json:"id"`
	AuthorName       string `json:"author_name"`
	AuthorId         string `json:"author_id"`
	AuthorChannelURL string `json:"author_channel_url"`
	Message          string `json:"message"`
	Likes            int    `json:"likes"`
	PublishedAt      string `json:"published_at"`
	ReplyCount       int    `json:"reply_count"`
}

// Represents the JSON for one reply. Part of CommentOutbound.
type Reply struct {
	Type             string `json:"type"`
	Id               string `json:"id"`
	ParentId         string `json:"parent_id"`
	AuthorName       string `json:"author_name"`
	AuthorId         string `json:"author_id"`
	AuthorChannelURL string `json:"author_channel_url"`
	Message          string `json:"message"`
	Likes            int    `json:"likes"`
	PublishedAt      string `json:"published_at"`
}

// Represents the JSON sent by the Comment endpoint.
type CommentOutbound struct {
	QuotaUsage int           `json:"quota_usage"`
	VideoId    string        `json:"video_id"`
	Comments   []interface{} `json:"comments"`
}

// Represents the JSON for a filter query.
type Filter struct {
	CaseSensitive bool     `json:"case_sensitive"`
	MatchAny      bool     `json:"match_any"`
	Reductive     bool     `json:"reductive"`
	Users         []string `json:"users"`
	Content       []string `json:"content"`
}
