package yt_stats

import "time"

type YoutubeErrorRequest struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Errors  []struct {
			Reason string `json:"reason"`
		} `json:"errors"`
	} `json:"error"`
}

type Test struct {
	Kind     string `json:"kind"`
	Etag     string `json:"etag"`
	PageInfo struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []struct {
		Kind string `json:"kind"`
		Etag string `json:"etag"`
		Id   string `json:"id"`
	} `json:"items"`
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
