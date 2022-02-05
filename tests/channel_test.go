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
	if !reflect.DeepEqual(outbound, expected) {
		t.Errorf("function parsed struct incorrectly: expected %+v actually %+v", expected, outbound)
	}
}

func TestChannelHandlerSuccess(t *testing.T) {
	var response yt_stats.ChannelOutbound
	req, err := http.NewRequest("GET", fmt.Sprintf("/ytstats/v1/channel/?id=%s", ChannelId), nil)
	req.Header.Set("key", getTestKey(t))
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
	if response.Channels[0].Id != ChannelId {
		t.Error("handler returned wrong body, got back wrong channel id.")
	}
	if response.QuotaUsage != 1 {
		t.Errorf("handler returned wrong quota usage: expected 1 actually %d", response.QuotaUsage)
	}
}

func TestChannelHandlerInvalidKey(t *testing.T) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/ytstats/v1/channel/?id=%s", ChannelId), nil)
	req.Header.Set("key", "invalid")
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := yt_stats.ChannelHandler(getInputs())
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: expected %v actually %v", http.StatusBadRequest, status)
	}
	expected := fmt.Sprintf(`{"quota_usage":0,"status_code":%d,"status_message":"keyInvalid"}`, http.StatusBadRequest)
	if strings.Trim(rr.Body.String(), "\n") != expected {
		t.Errorf("handler returned wrong body: expected %v actually %v", expected, rr.Body.String())
	}
}

func TestChannelHandlerNoKey(t *testing.T) {
	keyMissing(t, yt_stats.ChannelHandler, fmt.Sprintf("/ytstats/v1/channel/?id=%s", ChannelId))
}

func TestChannelHandlerNoChannel(t *testing.T) {
	req, err := http.NewRequest("GET", "/ytstats/v1/channel/", nil)
	req.Header.Set("key", getTestKey(t))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := yt_stats.ChannelHandler(getInputs())
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

func TestChannelHandlerTooManyChannels(t *testing.T) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/ytstats/v1/channel/?id=%s", strings.Repeat(",", 50)), nil)
	req.Header.Set("key", getTestKey(t))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := yt_stats.ChannelHandler(getInputs())
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: expected %v actually %v", http.StatusBadRequest, status)
	}
	expected := fmt.Sprintf(`{"quota_usage":0,"status_code":%d,"status_message":"tooManyItems"}`, http.StatusBadRequest)
	if strings.Trim(rr.Body.String(), "\n") != expected {
		t.Errorf("handler returned wrong body: expected %v actually %v", expected, rr.Body.String())
	}
}

func TestChannelHandlerUnsupportedType(t *testing.T) {
	unsupportedRequestType(t, yt_stats.ChannelHandler, "/ytstats/v1/channel/", "PUT")
}
