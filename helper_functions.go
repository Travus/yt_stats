package yt_stats

import (
	"encoding/json"
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

func getKey(w http.ResponseWriter, r *http.Request) string {
	key := r.URL.Query().Get("key")
	if key == "" {
		sendStatusCode(w, http.StatusBadRequest, "keyMissing")
	}
	return key
}
