package yt_stats

import (
	"time"
)

// Inputs stores variables sent to the handlers, basically global variables.
type Inputs struct {
	StartTime         time.Time
	StatusCheck       string
	RepliesRoot       string
	CommentsRoot      string
	ChannelsRoot      string
	PlaylistsRoot     string
	PlaylistItemsRoot string
	VideosRoot        string
	StreamRoot        string
	ChatRoot          string
}

// YoutubeErrorInbound represents the JSON received from a YouTube error response.
type YoutubeErrorInbound struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Errors  []struct {
			Reason string `json:"reason"`
		} `json:"errors"`
	} `json:"error"`
}

// StatusCodeOutbound represents the JSON sent on errors.
type StatusCodeOutbound struct {
	QuotaUsage    int    `json:"quota_usage"`
	StatusCode    int    `json:"status_code"`
	StatusMessage string `json:"status_message"`
}

// StatusOutbound represents the JSON sent by the Status endpoint.
type StatusOutbound struct {
	QuotaUsage    int     `json:"quota_usage"`
	Version       string  `json:"version"`
	Uptime        float64 `json:"uptime"`
	YoutubeStatus struct {
		StatusCode    int    `json:"status_code"`
		StatusMessage string `json:"status_message"`
	} `json:"youtube_status"`
}

// ChannelInbound represents the JSON received from the YouTube Channels endpoint.
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

// Channel represets the JSON for one channel. Part of ChannelOutbound struct.
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

// ChannelOutbound represets the JSON sent by the Channel endpoint.
type ChannelOutbound struct {
	QuotaUsage int       `json:"quota_usage"`
	Channels   []Channel `json:"channels"`
}

// PlaylistInbound represents the JSON received from the YouTube Playlists endpoint.
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

// PlaylistItemsInbound represents the JSON received from the YouTube PlaylistsItems endpoint.
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

// VideoInbound represents the JSON received from the YouTube Videos endpoint.
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
			CommentCount string `json:"commentCount"`
		} `json:"statistics"`
	} `json:"items"`
}

// VideoStats represents the JSON for stats over a range of videos. Part of Playlist struct.
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
	MostCommentedVideo    string `json:"most_commented_video"`
	MostComments          int    `json:"most_comments"`
	LeastCommentedVideo   string `json:"least_commented_video"`
	LeastComments         int    `json:"least_comments"`
	AverageComments       int    `json:"average_comments"`
}

// Video represents the JSON for one video. Part of Playlist struct.
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
	CommentCount int    `json:"comment_count"`
}

// Playlist represents the JSON for one playlist. Part of PlaylistOutbound struct.
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

// PlaylistOutbound represents the JSON sent by the Playlist endpoint.
type PlaylistOutbound struct {
	QuotaUsage int        `json:"quota_usage"`
	Playlists  []Playlist `json:"playlists"`
}

// VideoOutbound represents the JSON sent by the Video endpoint.
type VideoOutbound struct {
	QuotaUsage int         `json:"quota_usage"`
	VideoStats *VideoStats `json:"video_stats,omitempty"`
	Videos     []Video     `json:"videos"`
}

// CommentsInbound represents the JSON received from the YouTube CommentThreads endpoint.
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

// RepliesInbound represents the JSON received from the YouTube Comments endpoint.
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

// Comment represents the JSON for one comment. Part of CommentOutbound.
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

// Reply represents the JSON for one reply. Part of CommentOutbound.
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

// CommentOutbound represents the JSON sent by the Comment endpoint.
type CommentOutbound struct {
	QuotaUsage int           `json:"quota_usage"`
	VideoId    string        `json:"video_id"`
	Comments   []interface{} `json:"comments"`
}

// Filter represents the JSON for a filter query.
type Filter struct {
	CaseSensitive bool     `json:"case_sensitive"`
	MatchAny      bool     `json:"match_any"`
	Reductive     bool     `json:"reductive"`
	Users         []string `json:"users"`
	Content       []string `json:"content"`
}

// StreamInbound represents the JSON received from the YouTube Video endpoint used for Streams endpoint.
type StreamInbound struct {
	Items []struct {
		Id                   string `json:"id"`
		LiveStreamingDetails struct {
			ScheduledStartTime string `json:"scheduledStartTime"`
			ActualStartTime    string `json:"actualStartTime"`
			ActualEndTime      string `json:"actualEndTime"`
			ConcurrentViewers  string `json:"concurrentViewers"`
			ActiveLiveChatId   string `json:"activeLiveChatId"`
		} `json:"liveStreamingDetails"`
	} `json:"items"`
}

// StreamOutbound represents the JSON sent by the Streams endpoint.
type StreamOutbound struct {
	QuotaUsage int           `json:"quota_usage"`
	Streams    []interface{} `json:"streams"`
}

// LiveStream represents the JSON for one ongoing stream, needs own struct to display non-omitted viewer count.
type LiveStream struct {
	Id                 string `json:"id"`
	Status             string `json:"status"`
	ScheduledStartTime string `json:"scheduled_start_time,omitempty"`
	StartTime          string `json:"start_time"`
	ConcurrentViewers  int    `json:"concurrent_viewers"`
	ChatId             string `json:"chat_id,omitempty"`
}

// Stream represents the JSON for one not ongoing stream.
type Stream struct {
	Id                 string `json:"id"`
	Status             string `json:"status"`
	ScheduledStartTime string `json:"scheduled_start_time,omitempty"`
	StartTime          string `json:"start_time,omitempty"`
	EndTime            string `json:"end_time,omitempty"`
}

// ChatInbound represents the JSON received from the YouTube liveChatMessages endpoint.
type ChatInbound struct {
	NextPageToken         string `json:"nextPageToken"`
	PollingIntervalMillis int    `json:"pollingIntervalMillis"`
	Items                 []struct {
		Id      string `json:"id"`
		Snippet struct {
			Type                  string `json:"type"`
			AuthorChannelId       string `json:"authorChannelId"`
			PublishedAt           string `json:"publishedAt"`
			DisplayMessage        string `json:"displayMessage"`
			MessageDeletedDetails struct {
				DeletedMessageId string `json:"deletedMessageId"`
			} `json:"messageDeletedDetails"`
			UserBannedDetails struct {
				BannedUserDetails struct {
					ChannelId   string `json:"channelId"`
					ChannelUrl  string `json:"channelUrl"`
					DisplayName string `json:"displayName"`
				} `json:"bannedUserDetails"`
				BanType            string `json:"banType"`
				BanDurationSeconds int    `json:"banDurationSeconds"`
			} `json:"userBannedDetails"`
			SuperChatDetails struct {
				AmountMicros string `json:"amountMicros"`
				Currency     string `json:"currency"`
				UserComment  string `json:"userComment"`
			} `json:"superChatDetails"`
			SuperStickerDetails struct {
				SuperStickerMetadata struct {
					StickerId string `json:"stickerId"`
					AltText   string `json:"altText"`
				} `json:"superStickerMetadata"`
				AmountMicros string `json:"amountMicros"`
				Currency     string `json:"currency"`
			} `json:"superStickerDetails"`
			NewSponsorDetails struct {
				MemberLevelName string `json:"memberLevelName"`
				IsUpgrade       bool   `json:"isUpgrade"`
			}
			MemberMilestoneChatDetails struct {
				MemberLevelName string `json:"memberLevelName"`
				MemberMonth     int    `json:"memberMonth"`
				UserComment     string `json:"userComment"`
			} `json:"memberMilestoneChatDetails"`
		} `json:"snippet"`
		AuthorDetails struct {
			ChannelId       string `json:"channelId"`
			ChannelUrl      string `json:"channelUrl"`
			DisplayName     string `json:"displayName"`
			IsVerified      bool   `json:"isVerified"`
			IsChatOwner     bool   `json:"isChatOwner"`
			IsChatSponsor   bool   `json:"isChatSponsor"`
			IsChatModerator bool   `json:"isChatModerator"`
		} `json:"authorDetails"`
	} `json:"items"`
}

// ChatOutbound represents the JSON sent by the Chat endpoint.
type ChatOutbound struct {
	QuotaUsage        int           `json:"quota_usage"`
	ChatId            string        `json:"chat_id"`
	NextPageToken     string        `json:"page_token"`
	SuggestedCooldown int           `json:"suggested_cooldown"`
	ChatEvents        []interface{} `json:"chat_events"`
}

// ChatUser represents the JSON for a user in chat. Part of chat events.
type ChatUser struct {
	AuthorName       string `json:"author_name"`
	AuthorId         string `json:"author_id"`
	AuthorChannelUrl string `json:"author_channel_url"`
	ChatOwner        bool   `json:"chat_owner,omitempty"`
	Moderator        bool   `json:"moderator,omitempty"`
	Member           bool   `json:"member,omitempty"`
	Verified         bool   `json:"verified,omitempty"`
}

// ChatEnded represents the JSON for the chat ending. Part of ChatOutbound.
type ChatEnded struct {
	Id          string `json:"id"`
	Type        string `json:"type"`
	PublishedAt string `json:"published_at"`
}

// ChatMessageDeleted represents the JSON for a chat message being deleted. Part of ChatOutbound.
type ChatMessageDeleted struct {
	Id             string   `json:"id"`
	Type           string   `json:"type"`
	PublishedAt    string   `json:"published_at"`
	DeletedMessage string   `json:"deleted_message"`
	DeletedBy      ChatUser `json:"deleted_by"`
}

// ChatNewMember represents the JSON for a new member or member level change. Part of ChatOutbound.
type ChatNewMember struct {
	Id          string   `json:"id"`
	Type        string   `json:"type"`
	PublishedAt string   `json:"published_at"`
	Message     string   `json:"message"`
	Level       string   `json:"level"`
	Upgrade     bool     `json:"upgrade"`
	NewMember   ChatUser `json:"new_member"`
}

// ChatMemberMilestone represents the JSON for a member announcing membership renewal. Part of ChatOutbound.
type ChatMemberMilestone struct {
	Id          string   `json:"id"`
	Type        string   `json:"type"`
	PublishedAt string   `json:"published_at"`
	Message     string   `json:"message"`
	UserComment string   `json:"user_comment"`
	Level       string   `json:"level"`
	Months      int      `json:"months"`
	Member      ChatUser `json:"member"`
}

// ChatMemberOnlyModeEnded represents the JSON for a chat stopping member only mode. Part of ChatOutbound.
type ChatMemberOnlyModeEnded struct {
	Id          string   `json:"id"`
	Type        string   `json:"type"`
	PublishedAt string   `json:"published_at"`
	EndedBy     ChatUser `json:"ended_by"`
}

// ChatMemberOnlyModeStarted represents the JSON for a chat starting member only mode. Part of ChatOutbound.
type ChatMemberOnlyModeStarted struct {
	Id          string   `json:"id"`
	Type        string   `json:"type"`
	PublishedAt string   `json:"published_at"`
	StartedBy   ChatUser `json:"started_by"`
}

// ChatSuperChat represents the JSON for a super chat. Part of ChatOutbound.
type ChatSuperChat struct {
	Id          string   `json:"id"`
	Type        string   `json:"type"`
	PublishedAt string   `json:"published_at"`
	Message     string   `json:"message"`
	Amount      float64  `json:"amount"`
	Currency    string   `json:"currency"`
	SentBy      ChatUser `json:"sent_by"`
}

// ChatSuperSticker represents the JSON for a super sticker. Part of ChatOutbound.
type ChatSuperSticker struct {
	Id          string   `json:"id"`
	Type        string   `json:"type"`
	PublishedAt string   `json:"published_at"`
	Amount      float64  `json:"amount"`
	Currency    string   `json:"currency"`
	StickerId   string   `json:"sticker_id"`
	AltText     string   `json:"alt_text"`
	SentBy      ChatUser `json:"sent_by"`
}

// ChatMessage represents the JSON for chat message. Part of ChatOutbound.
type ChatMessage struct {
	Id          string   `json:"id"`
	Type        string   `json:"type"`
	PublishedAt string   `json:"published_at"`
	Message     string   `json:"message"`
	Author      ChatUser `json:"author"`
}

// ChatTombstone represents the JSON for a removed chat message. Part of ChatOutbound.
type ChatTombstone struct {
	Id          string `json:"id"`
	Type        string `json:"type"`
	PublishedAt string `json:"published_at"`
}

// ChatUserBanned represents the JSON for a chat user getting banned. Part of ChatOutbound.
type ChatUserBanned struct {
	Id          string   `json:"id"`
	Type        string   `json:"type"`
	PublishedAt string   `json:"published_at"`
	BanType     string   `json:"ban_type"`
	BanDuration int      `json:"ban_duration,omitempty"`
	BannedUser  ChatUser `json:"banned_user"`
	BannedBy    ChatUser `json:"banned_by"`
}

// ChatUnknownEvent represents the JSON for an unknown chat events the wrapper can't handle. Part of ChatOutbound.
type ChatUnknownEvent struct {
	Type  string      `json:"type"`
	Event interface{} `json:"event"`
}
