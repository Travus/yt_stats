package yt_stats

import "net/http"

func ChannelHandler(input Inputs) http.Handler {
	channel := func(w http.ResponseWriter, r *http.Request) {

	}
	return http.HandlerFunc(channel)
}
