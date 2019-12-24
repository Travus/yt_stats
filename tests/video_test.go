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

func getVideoIds(t *testing.T) string {
	var videoIds []string
	parseFile(t, "res/video_inbound.json", &videoIds)
	return strings.Join(videoIds, ",")
}

func TestVideoHandlerSuccess(t *testing.T) {
	var response yt_stats.VideoOutbound
	req, err := http.NewRequest("GET", fmt.Sprintf("/ytstats/v1/video/?key=%s&id=%s&stats=true",
		getKey(t), getVideoIds(t)), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := yt_stats.VideoHandler(getInputs())
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: expected %v actually %v", http.StatusOK, status)
	}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal("failed decoding response from endpoint")
	}
	if reflect.DeepEqual(response, yt_stats.VideoOutbound{}) {
		t.Error("function returned empty struct")
	}
	if response.Videos[0].Id != strings.Split(getVideoIds(t), ",")[0] {
		t.Error("handler returned wrong body, got back wrong video id")
	}
	if response.VideoStats == nil {
		t.Error("handler returned wrong body, got back no stats despite asking for them")
	}
	if response.QuotaUsage != 7 {
		t.Errorf("handler returned wrong quota usage: expected %d actually %d", 7, response.QuotaUsage)
	}
}

func TestVideoHandlerInvalidKey(t *testing.T) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/ytstats/v1/video/?key=invalid&id=%s", getVideoIds(t)), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := yt_stats.VideoHandler(getInputs())
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: expected %v actually %v", http.StatusBadRequest, status)
	}
	expected := fmt.Sprintf(`{"quota_usage":0,"status_code":%d,"status_message":"keyInvalid"}`, http.StatusBadRequest)
	if strings.Trim(rr.Body.String(), "\n") != expected {
		t.Errorf("handler returned wrong body: expected %v actually %v", expected, rr.Body.String())
	}
}

func TestVideoHandlerNoKey(t *testing.T) {
	keyMissing(t, yt_stats.VideoHandler, fmt.Sprintf("/ytstats/v1/video/?id=%s", getVideoIds(t)))
}

func TestVideoHandlerInvalidFlag(t *testing.T) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/ytstats/v1/playlist/?key=%s&id=%s&stats=invalid",
		getKey(t), getVideoIds(t)), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := yt_stats.VideoHandler(getInputs())
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: expected %v actually %v", http.StatusBadRequest, status)
	}
	expected := fmt.Sprintf(`{"quota_usage":0,"status_code":%d,"status_message":"flagInvalid"}`, http.StatusBadRequest)
	if strings.Trim(rr.Body.String(), "\n") != expected {
		t.Errorf("handler returned wrong body: expected %v actually %v", expected, rr.Body.String())
	}
}

func TestVideoHandlerNoVideo(t *testing.T) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/ytstats/v1/video/?key=%s", getKey(t)), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := yt_stats.VideoHandler(getInputs())
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: expected %v actually %v", http.StatusBadRequest, status)
	}
	expected := fmt.Sprintf(`{"quota_usage":0,"status_code":%d,"status_message":"channelIdMissing"}`,
		http.StatusBadRequest)
	if strings.Trim(rr.Body.String(), "\n") != expected {
		t.Errorf("handler returned wrong body: expected %v actually %v", expected, rr.Body.String())
	}
}

func TestVideoHandlerTooManyVideos(t *testing.T) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/ytstats/v1/video/?key=%s&id=%s",
		getKey(t), strings.Repeat(",", 50)), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := yt_stats.VideoHandler(getInputs())
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: expected %v actually %v", http.StatusBadRequest, status)
	}
	expected := fmt.Sprintf(`{"quota_usage":0,"status_code":%d,"status_message":"tooManyItems"}`, http.StatusBadRequest)
	if strings.Trim(rr.Body.String(), "\n") != expected {
		t.Errorf("handler returned wrong body: expected %v actually %v", expected, rr.Body.String())
	}
}

func TestVideoHandlerUnsupportedType(t *testing.T) {
	unsupportedRequestType(t, yt_stats.VideoHandler, "/ytstats/v1/video/", "PUT")
}
