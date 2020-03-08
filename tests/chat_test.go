package yt_stats_test

import (
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

const chatId = "Cg0KC01HNmd2Z0hGMEtJKicKGFVDcWFJWmNjNHBpWjkxMXl0ZWlQTFVXURILTUc2Z3ZnSEYwS0k"

// Required to convert sample data into the state it would otherwise be.
func fromFileFixerChat(t *testing.T, f string) yt_stats.ChatOutbound {
	var outbound yt_stats.ChatOutbound
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
	outbound.ChatEvents = make([]interface{}, len(inbound))
	for i, rawEntry := range inbound {
		if entry, ok := rawEntry.(map[string]interface{}); ok {
			if entry["type"] == "message" {
				event := yt_stats.ChatMessage{
					Id:          entry["id"].(string),
					Type:        "message",
					PublishedAt: entry["published_at"].(string),
					Message:     entry["message"].(string),
				}
				if author, authorOk := entry["author"].(map[string]interface{}); authorOk {
					event.Author = yt_stats.ChatUser{
						AuthorName:       author["author_name"].(string),
						AuthorId:         author["author_id"].(string),
						AuthorChannelUrl: author["author_channel_url"].(string),
						ChatOwner:        author["chat_owner"].(bool),
						Moderator:        author["moderator"].(bool),
						Sponsor:          author["sponsor"].(bool),
						Verified:         author["verified"].(bool),
					}
				}
				outbound.ChatEvents[i] = event
			} else if entry["type"] == "sponsor" {
				event := yt_stats.ChatNewSponsor{
					Id:          entry["id"].(string),
					Type:        "sponsor",
					PublishedAt: entry["published_at"].(string),
					Message:     entry["message"].(string),
				}
				if author, authorOk := entry["new_sponsor"].(map[string]interface{}); authorOk {
					event.NewSponsor = yt_stats.ChatUser{
						AuthorName:       author["author_name"].(string),
						AuthorId:         author["author_id"].(string),
						AuthorChannelUrl: author["author_channel_url"].(string),
						ChatOwner:        author["chat_owner"].(bool),
						Moderator:        author["moderator"].(bool),
						Sponsor:          author["sponsor"].(bool),
						Verified:         author["verified"].(bool),
					}
				}
				outbound.ChatEvents[i] = event
			} else if entry["type"] == "superchat" {
				event := yt_stats.ChatSuperChat{
					Id:          entry["id"].(string),
					Type:        "superchat",
					PublishedAt: entry["published_at"].(string),
					Message:     entry["message"].(string),
					Amount:      entry["amount"].(float64),
					Currency:    entry["currency"].(string),
				}
				if author, authorOk := entry["sent_by"].(map[string]interface{}); authorOk {
					event.SentBy = yt_stats.ChatUser{
						AuthorName:       author["author_name"].(string),
						AuthorId:         author["author_id"].(string),
						AuthorChannelUrl: author["author_channel_url"].(string),
						ChatOwner:        author["chat_owner"].(bool),
						Moderator:        author["moderator"].(bool),
						Sponsor:          author["sponsor"].(bool),
						Verified:         author["verified"].(bool),
					}
				}
				outbound.ChatEvents[i] = event
			} else {
				outbound.ChatEvents[i] = nil
			}
		}
	}
	return outbound
}

func TestChatParser(t *testing.T) {
	var inbound yt_stats.ChatInbound
	var expected yt_stats.ChatOutbound
	parseFile(t, "res/chat_inbound.json", &inbound)
	expected = fromFileFixerChat(t, "res/chat_outbound.json")
	outbound := yt_stats.ChatParser(inbound, chatId)
	if reflect.DeepEqual(outbound, yt_stats.ChatOutbound{}) {
		t.Error("function returned empty struct")
	}
	if !reflect.DeepEqual(outbound, expected) {
		t.Errorf("function parsed struct incorrectly: expected %+v actually %+v", expected, outbound)
	}
}

func TestChatHandlerInvalidKey(t *testing.T) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/ytstats/v1/chat/?key=invalid&id=%s", chatId), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := yt_stats.ChatHandler(getInputs())
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: expected %v actually %v", http.StatusBadRequest, status)
	}
	expected := fmt.Sprintf(`{"quota_usage":0,"status_code":%d,"status_message":"keyInvalid"}`, http.StatusBadRequest)
	if strings.Trim(rr.Body.String(), "\n") != expected {
		t.Errorf("handler returned wrong body: expected %v actually %v", expected, rr.Body.String())
	}
}

func TestChatHandlerNoKey(t *testing.T) {
	keyMissing(t, yt_stats.ChatHandler, fmt.Sprintf("/ytstats/v1/chat/?id=%s", chatId))
}

func TestChatHandlerNoChatId(t *testing.T) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/ytstats/v1/chat/?key=%s", getKey(t)), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := yt_stats.ChatHandler(getInputs())
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: expected %v actually %v", http.StatusBadRequest, status)
	}
	expected := fmt.Sprintf(`{"quota_usage":0,"status_code":%d,"status_message":"chatIdMissing"}`,
		http.StatusBadRequest)
	if strings.Trim(rr.Body.String(), "\n") != expected {
		t.Errorf("handler returned wrong body: expected %v actually %v", expected, rr.Body.String())
	}
}

func TestChatHandlerClosedChat(t *testing.T) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/ytstats/v1/chat/?key=%s&id=%s", getKey(t), chatId), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := yt_stats.ChatHandler(getInputs())
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: expected %v actually %v", http.StatusForbidden, status)
	}
	expected := fmt.Sprintf(`{"quota_usage":8,"status_code":%d,"status_message":"liveChatEnded"}`,
		http.StatusForbidden)
	if strings.Trim(rr.Body.String(), "\n") != expected {
		t.Errorf("handler returned wrong body: expected %v actually %v", expected, rr.Body.String())
	}
}

func TestChatHandlerUnsupportedType(t *testing.T) {
	unsupportedRequestType(t, yt_stats.ChatHandler, "/ytstats/v1/chat/", "PUT")
}
