package yt_stats_test

import (
	"fmt"
	"io/ioutil"
	"time"
	"yt_stats"
)

// Retrieves a valid API token from a 'token' file to run tests with.
func getKey() string {
	resp, err := ioutil.ReadFile("token")
	if err != nil {
		fmt.Print(err)
	}
	return string(resp)
}

// Gives inputs used by handlers.
func getInputs() yt_stats.Inputs {
	return yt_stats.Inputs{
		StartTime:    time.Now(),
		RepliesRoot:  "https://www.googleapis.com/youtube/v3/comments",
		CommentsRoot: "https://www.googleapis.com/youtube/v3/commentThreads",
		ChannelsRoot: "https://www.googleapis.com/youtube/v3/channels",
	}
}
