package yt_stats_test

import (
	"fmt"
	"reflect"
	"testing"
	"yt_stats"
)

func parseMockedPlaylist(t *testing.T, stats bool, videos bool) yt_stats.PlaylistOutbound {
	var inbound yt_stats.PlaylistInbound
	plVideos1 := make([]yt_stats.VideoInbound, 1)
	plVideos2 := make([]yt_stats.VideoInbound, 2)
	parseFile(t, "res/playlist_inbound.json", &inbound)
	parseFile(t, "res/video_inbound_1-1.json", &plVideos1[0])
	parseFile(t, "res/video_inbound_2-1.json", &plVideos2[0])
	parseFile(t, "res/video_inbound_2-2.json", &plVideos2[1])
	outbound := yt_stats.PlaylistTopLevelParser(inbound)
	outbound.Playlists = make([]yt_stats.Playlist, 2)
	yt_stats.VideoParser(plVideos1, &outbound.Playlists[0], stats, videos)
	yt_stats.VideoParser(plVideos2, &outbound.Playlists[1], stats, videos)
	return outbound
}

func TestPlaylistTopLevelParser(t *testing.T) {
	var inbound yt_stats.PlaylistInbound
	var expected yt_stats.PlaylistOutbound
	parseFile(t, "res/playlist_inbound.json", &inbound)
	parseFile(t, "res/top_level_playlist_outbound.json", expected)
	outbound := yt_stats.PlaylistTopLevelParser(inbound)
	if reflect.DeepEqual(outbound, yt_stats.Playlist{}) {
		t.Error("function returned empty struct")
	}
	if !reflect.DeepEqual(outbound, expected) {
		t.Errorf("function parsed struct incorrectly: expected %+v actually %+v", expected, outbound)
	}
}

func TestPlaylistItemsParser(t *testing.T) {
	inbound := make([]yt_stats.PlaylistItemsInbound, 2)
	var expected []string
	parseFile(t, "res/playlistitems_outbound.json", &expected)
	parseFile(t, "res/playlistitems_inbound_2-1.json", &inbound[0])
	parseFile(t, "res/playlistitems_inbound_2-2.json", &inbound[1])
	outbound := yt_stats.PlaylistItemsParser(inbound)
	if reflect.DeepEqual(outbound, []string{}) {
		t.Error("function returned empty data structure")
	}
	if !reflect.DeepEqual(outbound, expected) {
		t.Errorf("function parsed struct incorrectly: expected %+v actually %+v", expected, outbound)
	}
}

func TestVideoParser(t *testing.T) {
	videos1 := make([]yt_stats.VideoInbound, 1)
	videos2 := make([]yt_stats.VideoInbound, 2)
	var inbound yt_stats.PlaylistInbound
	var expected yt_stats.PlaylistOutbound
	outbound := yt_stats.PlaylistTopLevelParser(inbound)
	outbound.Playlists = make([]yt_stats.Playlist, 2)
	parseFile(t, "res/playlist_inbound.json", &outbound)
	parseFile(t, "res/video_inbound_1-1.json", &videos1[0])
	parseFile(t, "res/video_inbound_2-1.json", &videos2[0])
	parseFile(t, "res/video_inbound_2-2.json", &videos2[1])
	parseFile(t, "res/playlist_outbound.json", &expected)
	yt_stats.VideoParser(videos1, &outbound.Playlists[0], true, true)
	yt_stats.VideoParser(videos2, &outbound.Playlists[1], true, true)
	if !reflect.DeepEqual(outbound, expected) {
		t.Errorf("function parsed struct incorrectly: expected %+v actually %+v", expected, outbound)
	}
}

func TestFullPlaylistParsing(t *testing.T) {
	var expected yt_stats.PlaylistOutbound
	parseFile(t, "res/playlist_outbound.json", &expected)
	expected.Playlists[0].Videos = []yt_stats.Video{}
	expected.Playlists[1].VideoStats = yt_stats.VideoStats{}
	fmt.Printf("%+v", expected)
	outbound := parseMockedPlaylist(t, false, false)
	println("YO!")
	if reflect.DeepEqual(outbound, yt_stats.PlaylistOutbound{}) {
		t.Error("function returned empty data structure")
	}
	if !reflect.DeepEqual(outbound, expected) {
		t.Errorf("function parsed struct incorrectly: expected %+v actually %+v", expected, outbound)
	}

}

func TestPlaylistHandlerSuccess(t *testing.T) {

}

func TestPlaylistHandlerInvalidKey(t *testing.T) {

}

func TestPlaylistHandlerNoKey(t *testing.T) {

}

func TestPlaylistHandlerInvalidFlag(t *testing.T) {

}

func TestPlaylistHandlerNoPlaylist(t *testing.T) {

}

func TestPlaylistHandlerTooManyPlaylists(t *testing.T) {

}

func TestPlaylistHandlerUnsupportedType(t *testing.T) {

}
