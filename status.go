package yt_stats

import (
	"net/http"
)

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		http.Error(w, "Request type not supported.", http.StatusNotImplemented)
		break
	default:
		http.Error(w, "Request type not supported.", http.StatusNotImplemented)
	}
}
