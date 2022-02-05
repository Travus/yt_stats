package yt_stats_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"yt_stats"
)

const videoId = "zqfZs3Z7vy8"

// Required to convert sample data into the state it would otherwise be.
func fromFileFixerComments(t *testing.T, f string) []interface{} {
	var inbound []interface{}
	read, err := os.Open(f)
	if err != nil {
		t.Fatal(err)
	}
	err = json.NewDecoder(read).Decode(&inbound)
	if err != nil {
		t.Fatal(err)
	}
	err = read.Close()
	if err != nil {
		t.Fatal(err)
	}
	outbound := make([]interface{}, len(inbound))
	for i, rawEntry := range inbound {
		if entry, ok := rawEntry.(map[string]interface{}); ok {
			if entry["type"] == "comment" {
				outbound[i] = yt_stats.Comment{
					Type:             entry["type"].(string),
					Id:               entry["id"].(string),
					AuthorName:       entry["author_name"].(string),
					AuthorId:         entry["author_id"].(string),
					AuthorChannelURL: entry["author_channel_url"].(string),
					Message:          entry["message"].(string),
					Likes:            int(entry["likes"].(float64)),
					PublishedAt:      entry["published_at"].(string),
					ReplyCount:       int(entry["reply_count"].(float64)),
				}
			} else if entry["type"] == "reply" {
				outbound[i] = yt_stats.Reply{
					Type:             entry["type"].(string),
					Id:               entry["id"].(string),
					ParentId:         entry["parent_id"].(string),
					AuthorName:       entry["author_name"].(string),
					AuthorId:         entry["author_id"].(string),
					AuthorChannelURL: entry["author_channel_url"].(string),
					Message:          entry["message"].(string),
					Likes:            int(entry["likes"].(float64)),
					PublishedAt:      entry["published_at"].(string),
				}
			} else {
				outbound[i] = nil
			}
		}
	}
	return outbound
}

func TestCommentsParser(t *testing.T) {
	var inbound yt_stats.CommentsInbound
	var expected []interface{}
	var outbound []interface{}
	var replies []string
	parseFile(t, "res/commentthreads_inbound.json", &inbound)
	expected = fromFileFixerComments(t, "res/commentthreads_outbound.json")
	yt_stats.CommentsParser(inbound, &outbound, &replies)
	yt_stats.SortComments(&outbound)
	if reflect.DeepEqual(outbound, []yt_stats.Comment{}) {
		t.Error("function returned empty struct")
	}
	if len(replies) != 0 {
		t.Errorf("function parsed struct incorrectly: expected 0 replies actually %d", len(replies))
	}
	if !reflect.DeepEqual(outbound, expected) {
		t.Errorf("function parsed struct incorrectly: expected %+v actually %+v", expected, outbound)
	}
}

func TestRepliesParser(t *testing.T) {
	var inbound yt_stats.RepliesInbound
	var expected []interface{}
	var outbound []interface{}
	parseFile(t, "res/comments_inbound.json", &inbound)
	expected = fromFileFixerComments(t, "res/comments_outbound.json")
	yt_stats.RepliesParser(inbound, &outbound)
	yt_stats.SortComments(&outbound)
	if reflect.DeepEqual(outbound, []yt_stats.Comment{}) {
		t.Error("function returned empty struct")
	}
	if !reflect.DeepEqual(outbound, expected) {
		t.Errorf("function parsed struct incorrectly: expected %+v actually %+v", expected, outbound)
	}
}

func TestCommentSearch(t *testing.T) {
	var testData []interface{}
	var results []interface{}
	var filters []yt_stats.Filter
	testData = fromFileFixerComments(t, "res/sample_comments.json")
	parseFile(t, "res/searches_1.json", &filters)
	worked, results := yt_stats.CommentFilter(filters, testData)
	if !worked {
		t.Fatal("Failed filtering comments on search 1.")
	}
	if len(results) != 1 {
		t.Errorf("function returned wrong result, should have 1 comment, got %d", len(results))
	}
	results = nil
	filters = nil
	testData = fromFileFixerComments(t, "res/sample_comments.json")
	parseFile(t, "res/searches_2.json", &filters)
	worked, results = yt_stats.CommentFilter(filters, testData)
	if !worked {
		t.Fatal("Failed filtering comments on search 2.")
	}
	if len(results) != 5 {
		t.Errorf("function returned wrong result, should have 5 comments, got %d", len(results))
	}
	results = nil
	filters = nil
	testData = fromFileFixerComments(t, "res/sample_comments.json")
	parseFile(t, "res/searches_3.json", &filters)
	worked, results = yt_stats.CommentFilter(filters, testData)
	if !worked {
		t.Fatal("Failed filtering comments on search 3.")
	}
	if len(results) != 5 {
		t.Errorf("function returned wrong result, should have 5 comment, got %d", len(results))
	}
	results = nil
	filters = nil
	testData = fromFileFixerComments(t, "res/sample_comments.json")
	parseFile(t, "res/searches_4.json", &filters)
	worked, results = yt_stats.CommentFilter(filters, testData)
	if !worked {
		t.Fatal("Failed filtering comments on search 4.")
	}
	if len(results) != 0 {
		t.Errorf("function returned wrong result, should have 0 comments, got %d", len(results))
	}
	results = nil
	filters = nil
	testData = fromFileFixerComments(t, "res/sample_comments.json")
	parseFile(t, "res/searches_5.json", &filters)
	worked, results = yt_stats.CommentFilter(filters, testData)
	if !worked {
		t.Fatal("Failed filtering comments on search 5.")
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
		t.Errorf("handler returned wrong body: expected videoId %s, received videoId %s", videoId, response.VideoId)
	}
	if !(len(response.Comments) >= 48) {
		t.Errorf("handler returned wrong body: expected more than 48 results, received only %d", len(response.Comments))
	}
	if response.QuotaUsage < 1 {
		t.Error("handler returned low quota usage.")
	}
}

func TestCommentsHandlerSearchSuccess(t *testing.T) {
	var response yt_stats.CommentOutbound
	var filter []yt_stats.Filter
	parseFile(t, "res/search.json", &filter)
	body, err := json.Marshal(filter)
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
		t.Errorf("handler returned wrong body: expected videoId %s, received videoId %s", videoId, response.VideoId)
	}
	if len(response.Comments) > 20 || len(response.Comments) == 0 {
		t.Errorf("handler returned wrong body: got wrong number of comments, got %d", len(response.Comments))
	}
	if response.QuotaUsage < 1 {
		t.Error("handler returned low quota usage.")
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
	expected := fmt.Sprintf(`{"quota_usage":0,"status_code":%d,"status_message":"keyInvalid"}`, http.StatusBadRequest)
	if strings.Trim(rr.Body.String(), "\n") != expected {
		t.Errorf("handler returned wrong body: expected %v actually %v", expected, rr.Body.String())
	}
}

func TestCommentsHandlerInvalidSearch(t *testing.T) {
	body := "invalid"
	req, err := http.NewRequest("GET", fmt.Sprintf("/ytstats/v1/comments/?key=invalid&id=%s", videoId),
		strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := yt_stats.CommentsHandler(getInputs())
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: expected %v actually %v", http.StatusBadRequest, status)
	}
	expected := fmt.Sprintf(`{"quota_usage":0,"status_code":%d,"status_message":"searchBodyInvalid"}`,
		http.StatusBadRequest)
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
	expected := fmt.Sprintf(`{"quota_usage":0,"status_code":%d,"status_message":"videoIdMissing"}`,
		http.StatusBadRequest)
	if strings.Trim(rr.Body.String(), "\n") != expected {
		t.Errorf("handler returned wrong body: expected %v actually %v", expected, rr.Body.String())
	}
}

func TestCommentsHandlerUnsupportedType(t *testing.T) {
	unsupportedRequestType(t, yt_stats.CommentsHandler, "/ytstats/v1/comments/", "PUT")
}
