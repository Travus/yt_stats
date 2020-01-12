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

const streamId = "Qj9Ck1c3Zg0"

func TestStreamParser(t *testing.T) {
	var inbound yt_stats.StreamInbound
	var expected yt_stats.StreamOutbound
	parseFile(t, "res/stream_inbound.json", &inbound)
	parseFile(t, "res/stream_outbound.json", &expected)
	outbound := yt_stats.StreamParser(inbound)
	if reflect.DeepEqual(outbound, yt_stats.StreamOutbound{}) {
		t.Error("function returned empty struct")
	}
	if !reflect.DeepEqual(outbound, expected) {
		t.Errorf("function parsed struct incorrectly: expected %+v actually %+v", expected, outbound)
	}
}

func TestStreamHandlerInvalidKey(t *testing.T) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/ytstats/v1/stream/?key=invalid&id=%s", streamId), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := yt_stats.StreamHandler(getInputs())
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: expected %v actually %v", http.StatusBadRequest, status)
	}
	expected := fmt.Sprintf(`{"quota_usage":0,"status_code":%d,"status_message":"keyInvalid"}`, http.StatusBadRequest)
	if strings.Trim(rr.Body.String(), "\n") != expected {
		t.Errorf("handler returned wrong body: expected %v actually %v", expected, rr.Body.String())
	}
}

func TestStreamHandlerNoKey(t *testing.T) {
	keyMissing(t, yt_stats.StreamHandler, fmt.Sprintf("/ytstats/v1/stream/?id=%s", streamId))
}

func TestStreamHandlerNoVideo(t *testing.T) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/ytstats/v1/stream/?key=%s", getKey(t)), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := yt_stats.StreamHandler(getInputs())
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: expected %v actually %v", http.StatusBadRequest, status)
	}
	expected := fmt.Sprintf(`{"quota_usage":0,"status_code":%d,"status_message":"videoIdMissing"}`,
		http.StatusBadRequest)
	if strings.Trim(rr.Body.String(), "\n") != expected {
		t.Errorf("handler returned wrong body: expected %v actually %v", expected, rr.Body.String())
	}
}

func TestStreamHandlerTooManyVideos(t *testing.T) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/ytstats/v1/stream/?key=%s&id=%s",
		getKey(t), strings.Repeat(",", 50)), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := yt_stats.StreamHandler(getInputs())
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: expected %v actually %v", http.StatusBadRequest, status)
	}
	expected := fmt.Sprintf(`{"quota_usage":0,"status_code":%d,"status_message":"tooManyItems"}`, http.StatusBadRequest)
	if strings.Trim(rr.Body.String(), "\n") != expected {
		t.Errorf("handler returned wrong body: expected %v actually %v", expected, rr.Body.String())
	}
}

func TestStreamHandlerSuccess(t *testing.T) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/ytstats/v1/stream/?key=%s&id=%s", getKey(t), streamId), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := yt_stats.StreamHandler(getInputs())
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: expected %v actually %v", http.StatusOK, status)
	}
	var response yt_stats.StreamOutbound
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	if response.QuotaUsage != 3 || len(response.Streams) != 1 {
		t.Errorf("handler returned wrong body: expected 3 quota and 1 stream actually %d quota and %d streams",
			response.QuotaUsage, len(response.Streams))
	}
}

func TestStreamHandlerUnsupportedType(t *testing.T) {
	unsupportedRequestType(t, yt_stats.VideoHandler, "/ytstats/v1/stream/", "PUT")
}