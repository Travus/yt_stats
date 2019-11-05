package yt_stats

import (
	"encoding/json"
	"io"
	"net/http"
)

func ErrorParser(r io.Reader) StatusCodeResponse {
	var errorCode YoutubeErrorRequest
	err := json.NewDecoder(r).Decode(&errorCode)
	if err != nil {
		return StatusCodeResponse{StatusCode: http.StatusInternalServerError,
			                      StatusMessage: "yt_stats API failed to query YouTube"}
	}
	if errorCode.Error.Code != 0 {
		return StatusCodeResponse{StatusCode: errorCode.Error.Code, StatusMessage: errorCode.Error.Message}
	}
	return StatusCodeResponse{StatusCode: 200, StatusMessage: "OK"}
}  // ToDo: Consider using reason as message instead.
