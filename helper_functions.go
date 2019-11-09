package yt_stats

import (
	"encoding/json"
	"log"
	"net/http"
)

func unsupportedRequestType(w http.ResponseWriter) {
	response := StatusCodeOutbound{
		StatusCode:    http.StatusNotImplemented,
		StatusMessage: "requestTypeNotSupported",
	}
	w.WriteHeader(response.StatusCode)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Failed to respond to status endpoint.")
	}
}
