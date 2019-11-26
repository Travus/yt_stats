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

	inputs := yt_stats.Inputs{
		StartTime:             time.Now(),
		RepliesRoot:           "https://www.googleapis.com/youtube/v3/comments?part=snippet&maxResults=100&textFormat=plainText",
		CommentsRoot:          "https://www.googleapis.com/youtube/v3/commentThreads?part=snippet,replies&maxResults=100&textFormat=plainText",
		ChannelsRoot:          "https://www.googleapis.com/youtube/v3/channels?part=id,snippet,contentDetails,statistics&maxResults=50",
		PlaylistsRoot:         "https://www.googleapis.com/youtube/v3/playlists?part=snippet,contentDetails&maxResults=50",
		PlaylistItemsRootRoot: "https://www.googleapis.com/youtube/v3/playlistItems?part=snippet&maxResults=50",
		VideosRoot:            "https://www.googleapis.com/youtube/v3/videos?part=snippet,contentDetails,statistics&maxResults=50",
	}

	// Set port to 8080 and start handlers.
	port := "8080"
	http.HandleFunc("/ytstats/v1/", defaultHandler)
	http.Handle("/ytstats/v1/status/", yt_stats.StatusHandler(inputs))
	http.Handle("/ytstats/v1/channel/", yt_stats.ChannelHandler(inputs))
	http.Handle("/ytstats/v1/playlist/", yt_stats.PlaylistHandler(inputs))
	http.Handle("/ytstats/v1/video/", yt_stats.VideoHandler(inputs))

	// Serve REST API.
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}
