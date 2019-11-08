package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"yt_stats"
)

// Handles api root endpoint. /ytstats/v1/
// Just delivers basic info about the API's purpose.
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		_, err := fmt.Fprintf(w, "This API provides an easy interface to filter youtube comments by author "+
			"and comments.\nIt also provides an interface to get youtube channel stats.")
		if err != nil {
			log.Printf("Something went wrong writing REST API info.")
		}
		break
	default:
		http.Error(w, "Request type not supported.", http.StatusNotImplemented)
	}
}

func main() {

	yt_stats.StartTime = time.Now()
	yt_stats.CommentRoot = "https://www.googleapis.com/youtube/v3/comment"
	yt_stats.CommentsRoot = "https://www.googleapis.com/youtube/v3/commentThreads"
	yt_stats.ChannelsRoot = "https://www.googleapis.com/youtube/v3/channels"

	// Set port to 8080 and start handlers.
	port := "8080"
	http.HandleFunc("/ytstats/v1/", defaultHandler)
	http.HandleFunc("/ytstats/v1/status/", yt_stats.StatusHandler)

	// Serve REST API.
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}
