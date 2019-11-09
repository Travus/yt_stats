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
	req, err := http.NewRequest("GET", fmt.Sprintf("/ytstats/v1/status/?key=%s", getKey()), nil)
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
	expected := `{"version":"v1","uptime":2,"youtube_status":{"status_code":200,"status_message":"OK"}}`
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
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: expected %v actually %v", http.StatusOK, status)
	}
	expected := `{"version":"v1","uptime":1,"youtube_status":{"status_code":400,"status_message":"keyInvalid"}}`
	if strings.Trim(rr.Body.String(), "\n") != expected {
		t.Errorf("handler returned wrong body: expected %v actually %v", expected, rr.Body.String())
	}
}

func TestStatusHandlerNoKey(t *testing.T) {
	req, err := http.NewRequest("GET", "/ytstats/v1/status/", nil)
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
	expected := `{"status_code":400,"status_message":"keyMissing"}`
	if strings.Trim(rr.Body.String(), "\n") != expected {
		t.Errorf("handler returned wrong body: expected %v actually %v", expected, rr.Body.String())
	}
}