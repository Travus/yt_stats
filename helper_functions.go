package yt_stats

import (
	"encoding/json"
	"log"
	"net/http"
)

func unsupportedRequestType(w http.ResponseWriter) {
	response := StatusCodeOutbound{
		StatusCode:    http.StatusMethodNotAllowed,
		StatusMessage: "methodNotSupported",
	}
	w.WriteHeader(response.StatusCode)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Failed to respond to status endpoint.")
	}
}

func getKey(w http.ResponseWriter, r *http.Request) string {
	key := r.URL.Query().Get("key")
	if key == "" {
		response := StatusCodeOutbound{
			StatusCode:    http.StatusBadRequest,
			StatusMessage: "keyMissing",
		}
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Println("Failed to respond to status endpoint.")
		}
	}
	return key
}
