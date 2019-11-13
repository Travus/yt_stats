package yt_stats_test

import (
	"reflect"
	"testing"
	"yt_stats"
)

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
	expected := make([]string, 96)
	parseFile(t, "res/playlistItems_inbound_2-1.json", &inbound[0])
	parseFile(t, "res/playlistItems_inbound_2-2.json", &inbound[1])
	parseFile(t, "res/top_level_playlist_outbound.json", expected)
	outbound := yt_stats.PlaylistItemsParser(inbound)
	if reflect.DeepEqual(outbound, []string{}) {
		t.Error("function returned empty data structure.")
	}
	if !reflect.DeepEqual(outbound, expected) {
		t.Errorf("function parsed struct incorrectly: expected %+v actually %+v", expected, outbound)
	}
}

func TestVideoParser(t *testing.T) {
	inbound := make([]yt_stats.VideoInbound, 2)
	var outbound yt_stats.PlaylistOutbound
	var expected yt_stats.PlaylistOutbound
	parseFile(t, "res/video_inbound_2-1.json", &inbound[0])
	parseFile(t, "res/video_inbound_2-2.json", &inbound[1])
	parseFile(t, "res/playlist_outbound.json", expected)
	yt_stats.VideoParser(inbound, &outbound.Playlists[1], true, true)
	if !reflect.DeepEqual(outbound, expected) {
		t.Errorf("function parsed struct incorrectly: expected %+v actually %+v", expected, outbound)
	}
}

func TestFullPlaylistParsing(t *testing.T) {

}

func TestFullPlaylistParsingStats(t *testing.T) {

}

func TestFullPlaylistParsingVideos(t *testing.T) {

}

func TestFullPlaylistParsingAll(t *testing.T) {

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
