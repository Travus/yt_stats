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

const ChannelId = "UCBR8-60-B28hp2BmDPdntcQ"

func TestChannelParser(t *testing.T) {
	var inbound yt_stats.ChannelInbound
	var expected yt_stats.ChannelOutbound
	parseFile(t, "res/channel_inbound.json", &inbound)
	parseFile(t, "res/channel_outbound.json", &expected)
	outbound := yt_stats.ChannelParser(inbound)
	if reflect.DeepEqual(outbound, yt_stats.ChannelOutbound{}) {
		t.Error("function returned empty struct")
	}
	if reflect.DeepEqual(outbound, expected) {
		t.Errorf("function parsed struct incorrectly: expected %v actually %v", expected, outbound)
	}
}

func TestChannelHandlerSuccess(t *testing.T) {
	var response yt_stats.ChannelOutbound
	var expected yt_stats.ChannelOutbound
	parseFile(t, "res/channel_outbound.json", &expected)
	req, err := http.NewRequest("GET", fmt.Sprintf("/ytstats/v1/channel/?key=%s&id=%s", getKey(t), ChannelId), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := yt_stats.ChannelHandler(getInputs())
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: expected %v actually %v", http.StatusOK, status)
	}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal("failed decoding response from endpoint")
	}
	if reflect.DeepEqual(response, expected) {
		t.Errorf("handler returned wrong body: expected %v actually %v", expected, response)
	}
}

func TestChannelHandlerInvalidKey(t *testing.T) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/ytstats/v1/channel/?key=invalid&id=%s", ChannelId), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := yt_stats.ChannelHandler(getInputs())
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: expected %v actually %v", http.StatusOK, status)
	}
	expected := `{"version":"v1","uptime":1,"youtube_status":{"status_code":400,"status_message":"keyInvalid"}}`
	if strings.Trim(rr.Body.String(), "\n") != expected {
		t.Errorf("handler returned wrong body: expected %v actually %v", expected, rr.Body.String())
	}
}

func TestChannelHandlerNoKey(t *testing.T) {
	keyMissing(t, yt_stats.ChannelHandler, fmt.Sprintf("/ytstats/v1/channel/?id=%s", ChannelId))
}

func TestChannelHandlerNoChannel(t *testing.T) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/ytstats/v1/channel/?key=%s", getKey(t)), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := yt_stats.ChannelHandler(getInputs())
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: expected %v actually %v", http.StatusBadRequest, status)
	}
	expected := `{"status_code":400,"status_message":"channelIdMissing"}`
	if strings.Trim(rr.Body.String(), "\n") != expected {
		t.Errorf("handler returned wrong body: expected %v actually %v", expected, rr.Body.String())
	}
}