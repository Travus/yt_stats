package yt_stats

import (
	"encoding/json"
	duration "github.com/channelmeter/iso8601duration"
	"log"
	"net/http"
)

func sendStatusCode(w http.ResponseWriter, code int, msg string) {
	response := StatusCodeOutbound{
		StatusCode:    code,
		StatusMessage: msg,
	}
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Failed to respond to status endpoint.")
	}
}

func unsupportedRequestType(w http.ResponseWriter) {
	sendStatusCode(w, http.StatusMethodNotAllowed, "methodNotSupported")
}

func durationConverter(durStr string) (int, error) {
	dur, err := duration.FromString(durStr)
	if err != nil {
		return 0, err
	}
	return int(dur.ToDuration().Seconds()), nil
}
