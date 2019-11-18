package yt_stats_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"yt_stats"
)

const PlaylistIds = "PLbpi6ZahtOH7EUqtnmgJ3RFsej5UeFgxU%2CPLpjK416fmKwR-wFOaITVZ4Ktx2-mm2qp7"

func parseMockedPlaylist(t *testing.T, stats bool, videos bool) yt_stats.PlaylistOutbound {
	var inbound yt_stats.PlaylistInbound
	plVideos1 := make([]yt_stats.VideoInbound, 1)
	plVideos2 := make([]yt_stats.VideoInbound, 2)
	parseFile(t, "res/playlist_inbound.json", &inbound)
	parseFile(t, "res/video_inbound_1-1.json", &plVideos1[0])
	parseFile(t, "res/video_inbound_2-1.json", &plVideos2[0])
	parseFile(t, "res/video_inbound_2-2.json", &plVideos2[1])
	outbound := yt_stats.PlaylistTopLevelParser(inbound)
	// outbound.Playlists = make([]yt_stats.Playlist, 2)
	err := yt_stats.VideoParser(plVideos1, &outbound.Playlists[0], stats, videos)
	if err != nil {
		t.Fatal(err)
	}
	err = yt_stats.VideoParser(plVideos2, &outbound.Playlists[1], stats, videos)
	if err != nil {
		t.Fatal(err)
	}
	return outbound
}

func TestPlaylistTopLevelParser(t *testing.T) {
	var inbound yt_stats.PlaylistInbound
	var expected yt_stats.PlaylistOutbound
	parseFile(t, "res/playlist_inbound.json", &inbound)
	parseFile(t, "res/top_level_playlist_outbound.json", &expected)
	outbound := yt_stats.PlaylistTopLevelParser(inbound)
	if reflect.DeepEqual(outbound, yt_stats.PlaylistOutbound{}) {
		t.Error("function returned empty struct")
	}
	if !reflect.DeepEqual(outbound, expected) {
		t.Errorf("function parsed struct incorrectly: expected %+v actually %+v", expected, outbound)
	}
}

func TestPlaylistItemsParser(t *testing.T) {
	inbound := make([]yt_stats.PlaylistItemsInbound, 2)
	var expected [][]string
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
	parseFile(t, "res/playlist_inbound.json", &inbound)
	parseFile(t, "res/video_inbound_1-1.json", &videos1[0])
	parseFile(t, "res/video_inbound_2-1.json", &videos2[0])
	parseFile(t, "res/video_inbound_2-2.json", &videos2[1])
	parseFile(t, "res/playlist_outbound.json", &expected)
	outbound := yt_stats.PlaylistTopLevelParser(inbound)
	err := yt_stats.VideoParser(videos1, &outbound.Playlists[0], true, true)
	if err != nil {
		t.Fatal(err)
	}
	err = yt_stats.VideoParser(videos2, &outbound.Playlists[1], true, true)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(outbound, expected) {
		t.Errorf("function parsed struct incorrectly: expected %+v actually %+v", expected, outbound)
	}
}

func TestFullPlaylistParsing(t *testing.T) {
	var expected yt_stats.PlaylistOutbound
	parseFile(t, "res/playlist_outbound.json", &expected)
	expected.Playlists[0].Videos = nil
	expected.Playlists[1].Videos = nil
	expected.Playlists[0].VideoStats = yt_stats.VideoStats{}
	expected.Playlists[1].VideoStats = yt_stats.VideoStats{}
	outbound := parseMockedPlaylist(t, false, false)
	if reflect.DeepEqual(outbound, yt_stats.PlaylistOutbound{}) {
		t.Error("function returned empty data structure")
	}
	if !reflect.DeepEqual(outbound, expected) {
		t.Errorf("function parsed struct incorrectly: expected %+v actually %+v", expected, outbound)
	}
	parseFile(t, "res/playlist_outbound.json", &expected)
	expected.Playlists[0].Videos = nil
	expected.Playlists[1].Videos = nil
	outbound = parseMockedPlaylist(t, true, false)
	if reflect.DeepEqual(outbound, yt_stats.PlaylistOutbound{}) {
		t.Error("function returned empty data structure")
	}
	if !reflect.DeepEqual(outbound, expected) {
		t.Errorf("function parsed struct incorrectly: expected %+v actually %+v", expected, outbound)
	}
	parseFile(t, "res/playlist_outbound.json", &expected)
	expected.Playlists[0].VideoStats = yt_stats.VideoStats{}
	expected.Playlists[1].VideoStats = yt_stats.VideoStats{}
	outbound = parseMockedPlaylist(t, false, true)
	if reflect.DeepEqual(outbound, yt_stats.PlaylistOutbound{}) {
		t.Error("function returned empty data structure")
	}
	if !reflect.DeepEqual(outbound, expected) {
		t.Errorf("function parsed struct incorrectly: expected %+v actually %+v", expected, outbound)
	}
	parseFile(t, "res/playlist_outbound.json", &expected)
	outbound = parseMockedPlaylist(t, true, true)
	if reflect.DeepEqual(outbound, yt_stats.PlaylistOutbound{}) {
		t.Error("function returned empty data structure")
	}
	if !reflect.DeepEqual(outbound, expected) {
		t.Errorf("function parsed struct incorrectly: expected %+v actually %+v", expected, outbound)
	}
}

func TestPlaylistHandlerSuccess(t *testing.T) {
	var response yt_stats.PlaylistOutbound
	req, err := http.NewRequest("GET", fmt.Sprintf("/ytstats/v1/playlist/?key=%s&id=%s&videos=%s&stats=%s",
		getKey(t), PlaylistIds, "true", "false"), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := yt_stats.PlaylistHandler(getInputs())
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: expected %v actually %v", http.StatusOK, status)
	}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal("failed decoding response from endpoint")
	}
	if len(response.Playlists) != 2 {
		t.Error("handler returned wrong body, got back wrong amount of playlists")
	}
	if response.Playlists[0].VideoStats.AvailableVideos != 0 {
		t.Error("handler returned wrong body, got back stats despite not asking for them")
	}
	if len(response.Playlists[0].Videos) == 0 {
		t.Error("handler returned wrong body, got back no videos despite asking for them")
	}
	if response.Playlists[0].Id + "%2C" + response.Playlists[1].Id == PlaylistIds {
		t.Error("handler returned wrong body, got back wrong playlist ids")
	}
}

func TestPlaylistHandlerInvalidKey(t *testing.T) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/ytstats/v1/playlist/?key=invalid&id=%s",
		PlaylistIds), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := yt_stats.PlaylistHandler(getInputs())
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: expected %v actually %v", http.StatusBadRequest, status)
	}
	expected := fmt.Sprintf(`{"status_code":%d,"status_message":"keyInvalid"}`, http.StatusBadRequest)
	if strings.Trim(rr.Body.String(), "\n") != expected {
		t.Errorf("handler returned wrong body: expected %v actually %v", expected, rr.Body.String())
	}
}

func TestPlaylistHandlerNoKey(t *testing.T) {
	keyMissing(t, yt_stats.PlaylistHandler, fmt.Sprintf("/ytstats/v1/channel/?id=%s", PlaylistIds))
}

func TestPlaylistHandlerInvalidFlag(t *testing.T) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/ytstats/v1/playlist/?key=%s&id=%s&videos=invalid",
		getKey(t), PlaylistIds), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := yt_stats.PlaylistHandler(getInputs())
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: expected %v actually %v", http.StatusBadRequest, status)
	}
	expected := fmt.Sprintf(`{"status_code":%d,"status_message":"flagInvalid"}`, http.StatusBadRequest)
	if strings.Trim(rr.Body.String(), "\n") != expected {
		t.Errorf("handler returned wrong body: expected %v actually %v", expected, rr.Body.String())
	}
}

func TestPlaylistHandlerNoPlaylist(t *testing.T) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/ytstats/v1/playlist/?key=%s", getKey(t)), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := yt_stats.PlaylistHandler(getInputs())
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: expected %v actually %v", http.StatusBadRequest, status)
	}
	expected := fmt.Sprintf(`{"status_code":%d,"status_message":"playlistIdMissing"}`, http.StatusBadRequest)
	if strings.Trim(rr.Body.String(), "\n") != expected {
		t.Errorf("handler returned wrong body: expected %v actually %v", expected, rr.Body.String())
	}
}

func TestPlaylistHandlerTooManyPlaylists(t *testing.T) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/ytstats/v1/playlist/?key=%s&id=%s",
		getKey(t), strings.Repeat(",", 50)), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := yt_stats.PlaylistHandler(getInputs())
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: expected %v actually %v", http.StatusBadRequest, status)
	}
	expected := fmt.Sprintf(`{"status_code":%d,"status_message":"tooManyItems"}`, http.StatusBadRequest)
	if strings.Trim(rr.Body.String(), "\n") != expected {
		t.Errorf("handler returned wrong body: expected %v actually %v", expected, rr.Body.String())
	}
}

func TestPlaylistHandlerUnsupportedType(t *testing.T) {
	unsupportedRequestType(t, yt_stats.PlaylistHandler, "/ytstats/v1/playlist/", "PUT")
}
