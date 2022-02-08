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
func defaultHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			http.ServeFile(w, r, "html/ytstats.v1.html")
			break
		default:
			http.Error(w, "Request type not supported.", http.StatusNotImplemented)
		}
	})
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
		log.Fatal(http.ListenAndServe(":8080", certManager.HTTPHandler(nil)))
	}()
	log.Fatal(server.ListenAndServeTLS("", ""))
}

func runInDev(mux *http.ServeMux) {
	// Serve REST API.
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}

func logIncoming(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s -> %s", r.Method, r.RemoteAddr, r.URL.Path)
		handler.ServeHTTP(w, r)
	})
}

func main() {

	// Set global values.
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)
	inputs := yt_stats.Inputs{
		StartTime:         time.Now(),
		StatusCheck:       "https://www.googleapis.com/youtube/v3/channels?part=id&id=UCBR8-60-B28hp2BmDPdntcQ",
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
	mux.Handle("/ytstats/v1/", logIncoming(defaultHandler()))
	mux.Handle("/ytstats/v1/status/", logIncoming(yt_stats.StatusHandler(inputs)))
	mux.Handle("/ytstats/v1/channel/", logIncoming(yt_stats.ChannelHandler(inputs)))
	mux.Handle("/ytstats/v1/playlist/", logIncoming(yt_stats.PlaylistHandler(inputs)))
	mux.Handle("/ytstats/v1/video/", logIncoming(yt_stats.VideoHandler(inputs)))
	mux.Handle("/ytstats/v1/comments/", logIncoming(yt_stats.CommentsHandler(inputs)))
	mux.Handle("/ytstats/v1/stream/", logIncoming(yt_stats.StreamHandler(inputs)))
	mux.Handle("/ytstats/v1/chat/", logIncoming(yt_stats.ChatHandler(inputs)))

	if os.Getenv("tls_address") != "" {
		log.Print("Running in production mode...")
		runInProduction(mux)
	} else {
		log.Print("Running in development mode...")
		runInDev(mux)
	}
}
