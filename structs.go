package yt_stats

import "time"

type YoutubeErrorRequest struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type StatusCodeResponse struct {
	StatusCode    int    `json:"status_code"`
	StatusMessage string `json:"status_message"`
}

type StatusResponse struct {
	Version       string             `json:"version"`
	Uptime        time.Duration      `json:"uptime"`
	YoutubeStatus StatusCodeResponse `json:"youtube_status"`
}
