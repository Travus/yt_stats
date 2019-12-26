package yt_stats

import (
	"encoding/json"
	duration "github.com/channelmeter/iso8601duration"
	"log"
	"net/http"
)

// Sends a status code response, used to report back errors.
func sendStatusCode(w http.ResponseWriter, quota int, code int, msg string) {
	response := StatusCodeOutbound{
		QuotaUsage:    quota,
		StatusCode:    code,
		StatusMessage: msg,
	}
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Failed to send error response.")
	}
}

// Sends a generic status code for unsupported methods.
func unsupportedRequestType(w http.ResponseWriter) {
	sendStatusCode(w, 0, http.StatusMethodNotAllowed, "methodNotSupported")
}

// Converters ISO 8601 duration as string to amount of seconds as int.
func durationConverter(durStr string) (int, error) {
	dur, err := duration.FromString(durStr)
	if err != nil {
		return 0, err
	}
	return int(dur.ToDuration().Seconds()), nil
}
