package yt_stats

type YoutubeErrorRequest struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Errors  []struct {
			Reason string `json:"reason"`
		} `json:"errors"`
	} `json:"error"`
}

type StatusCodeResponse struct {
	StatusCode    int    `json:"status_code"`
	StatusMessage string `json:"status_message"`
}

type StatusResponse struct {
	Version       string             `json:"version"`
	Uptime        float64            `json:"uptime"`
	YoutubeStatus StatusCodeResponse `json:"youtube_status"`
}
