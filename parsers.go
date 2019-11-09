package yt_stats

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func ErrorParser(r io.Reader, s interface{}) StatusCodeOutbound {
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)
	var errorCode YoutubeErrorInbound
	err := json.NewDecoder(tee).Decode(&errorCode)
	if err != nil {
		return StatusCodeOutbound{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: "failedToQueryYouTubeAPI"}
	}
	if errorCode.Error.Code != 0 {
		if errorCode.Error.Errors[0].Reason == "" {
			return StatusCodeOutbound{StatusCode: errorCode.Error.Code, StatusMessage: errorCode.Error.Message}
		}
		return StatusCodeOutbound{StatusCode: errorCode.Error.Code, StatusMessage: errorCode.Error.Errors[0].Reason}
	}
	if s != nil {
		err := json.NewDecoder(&buf).Decode(&s)
		if err != nil {
			return StatusCodeOutbound{StatusCode: http.StatusInternalServerError,
				StatusMessage: "failedToQueryYouTubeAPI"}
		}
	}
	return StatusCodeOutbound{StatusCode: 200, StatusMessage: "OK"}
}
