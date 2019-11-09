package yt_stats

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// This function parses a youtube response into a given struct, and returns a status code struct with OK or response.
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
		err := json.NewDecoder(&buf).Decode(&s)
		if err != nil {
			return StatusCodeOutbound{
				StatusCode:    http.StatusInternalServerError,
				StatusMessage: "failedToQueryYouTubeAPI",
			}
		}
	}
	return StatusCodeOutbound{StatusCode: 200, StatusMessage: "OK"}
}

// This function parses a inbound channel struct to a outbound channel struct.
func ChannelParser(inbound ChannelInbound) ChannelOutbound {
	return ChannelOutbound{}
}
