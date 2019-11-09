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
