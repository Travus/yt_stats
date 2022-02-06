package main

import (
	"crypto/tls"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
	"os"
	"time"
	"yt_stats"
)

// Handler for api root endpoint. /ytstats/v1/
// Provides basic info about the API's purpose.
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeFile(w, r, "html/ytstats.v1.html")
		break
	default:
		http.Error(w, "Request type not supported.", http.StatusNotImplemented)
	}
}

func runInProduction(mux *http.ServeMux) {
	// Setup Automated Certificate Management Environment (ACME)
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache("cert-cache"),
		HostPolicy: autocert.HostWhitelist(os.Getenv("tls_address")),
	}

	// Serve REST API.
	server := &http.Server{
		Addr:    ":8081",
		Handler: mux,
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}
	go func() {
		err := http.ListenAndServe(":8080", certManager.HTTPHandler(nil))
		if err != nil {
			log.Fatal(err.Error())
		}
	}()
	err := server.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatal(err.Error())
	}
}

func runInDev(mux *http.ServeMux) {
	// Serve REST API.
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {

	// Set global values.
	inputs := yt_stats.Inputs{
		StartTime:         time.Now(),
		RepliesRoot:       "https://www.googleapis.com/youtube/v3/comments?part=snippet&maxResults=100&textFormat=plainText",
		CommentsRoot:      "https://www.googleapis.com/youtube/v3/commentThreads?part=snippet,replies&maxResults=100&textFormat=plainText",
		ChannelsRoot:      "https://www.googleapis.com/youtube/v3/channels?part=id,snippet,contentDetails,statistics&maxResults=50",
		PlaylistsRoot:     "https://www.googleapis.com/youtube/v3/playlists?part=snippet,contentDetails&maxResults=50",
		PlaylistItemsRoot: "https://www.googleapis.com/youtube/v3/playlistItems?part=snippet&maxResults=50",
		VideosRoot:        "https://www.googleapis.com/youtube/v3/videos?part=snippet,contentDetails,statistics&maxResults=50",
		StreamRoot:        "https://www.googleapis.com/youtube/v3/videos?part=id,liveStreamingDetails&maxResults=50",
		ChatRoot:          "https://www.googleapis.com/youtube/v3/liveChat/messages?part=id,snippet,authorDetails&maxResults=2000",
	}

	// Setup handlers.
	mux := http.NewServeMux()
	mux.HandleFunc("/ytstats/v1/", defaultHandler)
	mux.Handle("/ytstats/v1/status/", yt_stats.StatusHandler(inputs))
	mux.Handle("/ytstats/v1/channel/", yt_stats.ChannelHandler(inputs))
	mux.Handle("/ytstats/v1/playlist/", yt_stats.PlaylistHandler(inputs))
	mux.Handle("/ytstats/v1/video/", yt_stats.VideoHandler(inputs))
	mux.Handle("/ytstats/v1/comments/", yt_stats.CommentsHandler(inputs))
	mux.Handle("/ytstats/v1/stream/", yt_stats.StreamHandler(inputs))
	mux.Handle("/ytstats/v1/chat/", yt_stats.ChatHandler(inputs))

	if os.Getenv("tls_address") != "" {
		runInProduction(mux)
	} else {
		runInDev(mux)
	}
}
