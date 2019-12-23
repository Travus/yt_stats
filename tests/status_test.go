package yt_stats_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"yt_stats"
)

func TestStatusHandlerValidKey(t *testing.T) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/ytstats/v1/status/?key=%s", getKey(t)), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := yt_stats.StatusHandler(getInputs())
	time.Sleep(2 * time.Second)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: expected %v actually %v", http.StatusOK, status)
	}
	expected := fmt.Sprintf(`{"version":"v1","uptime":2,"youtube_status":{"status_code":%d,"status_message":"OK"}}`,
		http.StatusOK)
	if strings.Trim(rr.Body.String(), "\n") != expected {
		t.Errorf("handler returned wrong body: expected %v actually %v", expected, rr.Body.String())
	}
}

func TestStatusHandlerInvalidKey(t *testing.T) {
	req, err := http.NewRequest("GET", "/ytstats/v1/status/?key=invalid", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := yt_stats.StatusHandler(getInputs())
	time.Sleep(1 * time.Second)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: expected %v actually %v", http.StatusBadRequest, status)
	}
	expected := `{"quota_usage":1,"version":"v1","uptime":1,"youtube_status":` +
		fmt.Sprintf(`{"status_code":%d,"status_message":"keyInvalid"}}`, http.StatusBadRequest)
	if strings.Trim(rr.Body.String(), "\n") != expected {
		t.Errorf("handler returned wrong body: expected %v actually %v", expected, rr.Body.String())
	}
}

func TestStatusHandlerNoKey(t *testing.T) {
	keyMissing(t, yt_stats.StatusHandler, "/ytstats/v1/status/")
}

func TestStatusHandlerUnsupportedType(t *testing.T) {
	unsupportedRequestType(t, yt_stats.StatusHandler, "/ytstats/v1/status/", "PUT")
}
