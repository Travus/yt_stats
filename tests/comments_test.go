package yt_stats_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"yt_stats"
)

const videoId = "zqfZs3Z7vy8"

func TestCommentsParser(t *testing.T) {
	var inbound yt_stats.CommentsInbound
	var expected []yt_stats.Comment
	var outbound []yt_stats.Comment
	parseFile(t, "res/commentthreads_inbound.json", &inbound)
	parseFile(t, "res/commentthreads_outbound.json", &expected)
	err := yt_stats.CommentsParser(inbound, &outbound)
	if err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(outbound, []yt_stats.Comment{}) {
		t.Error("function returned empty struct")
	}
	if !reflect.DeepEqual(outbound, expected) {
		t.Errorf("function parsed struct incorrectly: expected %+v actually %+v", expected, outbound)
	}
}

func TestRepliesParser(t *testing.T) {
	var inbound yt_stats.RepliesInbound
	var expected []yt_stats.Comment
	var outbound []yt_stats.Comment
	parseFile(t, "res/comments_inbound.json", &inbound)
	parseFile(t, "res/comments_outbound.json", &expected)
	err := yt_stats.RepliesParser(inbound, &outbound)
	if err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(outbound, []yt_stats.Comment{}) {
		t.Error("function returned empty struct")
	}
	if !reflect.DeepEqual(outbound, expected) {
		t.Errorf("function parsed struct incorrectly: expected %+v actually %+v", expected, outbound)
	}
}

func TestCommentSearch(t *testing.T) {
	var testData []yt_stats.Comment
	var results []yt_stats.Comment
	var searches []yt_stats.Search
	parseFile(t, "res/sample_comments.json", &testData)
	parseFile(t, "res/searches_1.json", &searches)
	err := yt_stats.CommentSearch(searches, &results, &testData)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Errorf("function returned wrong result, should have 1 comment, got %d", len(results))
	}
	results = nil
	parseFile(t, "res/sample_comments.json", &testData)
	parseFile(t, "res/searches_2.json", &searches)
	err = yt_stats.CommentSearch(searches, &results, &testData)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 5 {
		t.Errorf("function returned wrong result, should have 5 comments, got %d", len(results))
	}
	results = nil
	parseFile(t, "res/sample_comments.json", &testData)
	parseFile(t, "res/searches_3.json", &searches)
	err = yt_stats.CommentSearch(searches, &results, &testData)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Errorf("function returned wrong result, should have 1 comment, got %d", len(results))
	}
	results = nil
	parseFile(t, "res/sample_comments.json", &testData)
	parseFile(t, "res/searches_4.json", &searches)
	err = yt_stats.CommentSearch(searches, &results, &testData)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Errorf("function returned wrong result, should have 0 comments, got %d", len(results))
	}
}

func TestCommentsHandlerSuccess(t *testing.T) {
	var response yt_stats.CommentOutbound
	req, err := http.NewRequest("GET", fmt.Sprintf("/ytstats/v1/comments/?key=%s&id=%s", getKey(t), videoId), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := yt_stats.CommentsHandler(getInputs())
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: expected %v actually %v", http.StatusOK, status)
	}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal("failed decoding response from endpoint")
	}
	if response.VideoId != videoId {
		t.Errorf("handler returned wrong body: expected videoId %s, received videoId %s",videoId, response.VideoId)
	}
	if !(len(response.Comments) >= 48) {
		t.Errorf("handler returned wrong body: expected more than 48 results, received only %d", len(response.Comments))
	}
}

func TestCommentsHandlerSearchSuccess(t *testing.T) {
	var response yt_stats.CommentOutbound
	var search []yt_stats.Search
	parseFile(t, "res/search.json", &search)
	body, err := json.Marshal(search)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("/ytstats/v1/comments/?key=%s&id=%s",
		getKey(t), videoId), bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := yt_stats.CommentsHandler(getInputs())
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: expected %v actually %v", http.StatusOK, status)
	}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal("failed decoding response from endpoint")
	}
	if response.VideoId != videoId {
		t.Errorf("handler returned wrong body: expected videoId %s, received videoId %s",videoId, response.VideoId)
	}
	if len(response.Comments) > 20 || len(response.Comments) == 0 {
		t.Errorf("handler returned wrong body: got wrong number of comments, got %d", len(response.Comments))
	}
}

func TestCommentsHandlerInvalidKey(t *testing.T) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/ytstats/v1/comments/?key=invalid&id=%s", videoId), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := yt_stats.CommentsHandler(getInputs())
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: expected %v actually %v", http.StatusBadRequest, status)
	}
	expected := fmt.Sprintf(`{"status_code":%d,"status_message":"keyInvalid"}`, http.StatusBadRequest)
	if strings.Trim(rr.Body.String(), "\n") != expected {
		t.Errorf("handler returned wrong body: expected %v actually %v", expected, rr.Body.String())
	}
}

func TestCommentsHandlerNoKey(t *testing.T) {
	keyMissing(t, yt_stats.CommentsHandler, fmt.Sprintf("/ytstats/v1/comments/?id=%s", videoId))
}


func TestCommentsHandlerNoVideo(t *testing.T) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/ytstats/v1/comments/?key=%s", getKey(t)), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := yt_stats.CommentsHandler(getInputs())
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: expected %v actually %v", http.StatusBadRequest, status)
	}
	expected := fmt.Sprintf(`{"status_code":%d,"status_message":"videoIdMissing"}`, http.StatusBadRequest)
	if strings.Trim(rr.Body.String(), "\n") != expected {
		t.Errorf("handler returned wrong body: expected %v actually %v", expected, rr.Body.String())
	}
}

func TestCommentsHandlerUnsupportedType(t *testing.T) {
	unsupportedRequestType(t, yt_stats.CommentsHandler, "/ytstats/v1/comments/", "PUT")
}

