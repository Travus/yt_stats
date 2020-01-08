package yt_stats_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
	"yt_stats"
)

// Retrieves a valid API token from a 'token' file to run tests with.
func getKey(t *testing.T) string {
	resp, err := ioutil.ReadFile("yt_token")
	if err != nil {
		t.Fatal(err)
	}
	return string(resp)
}

// Reads and parses json file into provided struct.
func parseFile(t *testing.T, f string, s interface{}) {
	read, err := os.Open(f)
	if err != nil {
		t.Fatal(err)
	}
	err = json.NewDecoder(read).Decode(&s)
	if err != nil {
		t.Fatal(err)
	}
	err = read.Close()
	if err != nil {
		t.Fatal(err)
	}
}

// Gives inputs used by handlers.
func getInputs() yt_stats.Inputs {
	return yt_stats.Inputs{
		StartTime:         time.Now(),
		RepliesRoot:       "https://www.googleapis.com/youtube/v3/comments?part=snippet&maxResults=100&textFormat=plainText",
		CommentsRoot:      "https://www.googleapis.com/youtube/v3/commentThreads?part=snippet,replies&maxResults=100&textFormat=plainText",
		ChannelsRoot:      "https://www.googleapis.com/youtube/v3/channels?part=id,snippet,contentDetails,statistics&maxResults=50",
		PlaylistsRoot:     "https://www.googleapis.com/youtube/v3/playlists?part=snippet,contentDetails&maxResults=50",
		PlaylistItemsRoot: "https://www.googleapis.com/youtube/v3/playlistItems?part=snippet&maxResults=50",
		VideosRoot:        "https://www.googleapis.com/youtube/v3/videos?part=snippet,contentDetails,statistics&maxResults=50",
	}
}

// Generalized test, testing output for URLs for missing key errors.
func keyMissing(t *testing.T, f func(inputs yt_stats.Inputs) http.Handler, url string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := f(getInputs())
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: expected %v actually %v", http.StatusBadRequest, status)
	}
	expected := fmt.Sprintf(`{"quota_usage":0,"status_code":%d,"status_message":"keyMissing"}`, http.StatusBadRequest)
	if strings.Trim(rr.Body.String(), "\n") != expected {
		t.Errorf("handler returned wrong body: expected %v actually %v", expected, rr.Body.String())
	}
}

// Generalized test, testing output for URLs for unsupported request type errors.
func unsupportedRequestType(t *testing.T, f func(inputs yt_stats.Inputs) http.Handler, url string, rType string) {
	req, err := http.NewRequest(rType, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := f(getInputs())
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: expected %v actually %v", http.StatusMethodNotAllowed, status)
	}
	expected := fmt.Sprintf(`{"quota_usage":0,"status_code":%d,"status_message":"methodNotSupported"}`,
		http.StatusMethodNotAllowed)
	if strings.Trim(rr.Body.String(), "\n") != expected {
		t.Errorf("handler returned wrong body: expected %v actually %v", expected, rr.Body.String())
	}
}
