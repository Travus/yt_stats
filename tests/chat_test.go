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
	var inbound interface{}
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
	if topLevel, ok := inbound.(map[string]interface{}); ok {
		outbound.ChatId = topLevel["chat_id"].(string)
		outbound.NextPageToken = topLevel["page_token"].(string)
		outbound.SuggestedCooldown = int(topLevel["suggested_cooldown"].(float64))
		outbound.ChatEvents = make([]interface{}, len(topLevel["chat_events"].([]interface{})))
		for i, rawEntry := range topLevel["chat_events"].([]interface{}) {
			if entry, entryOk := rawEntry.(map[string]interface{}); entryOk {
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
							Member:           author["member"].(bool),
							Verified:         author["verified"].(bool),
						}
					}
					outbound.ChatEvents[i] = event
				} else if entry["type"] == "new_member" {
					event := yt_stats.ChatNewMember{
						Id:          entry["id"].(string),
						Type:        "new_member",
						PublishedAt: entry["published_at"].(string),
						Message:     entry["message"].(string),
						Level:       entry["level"].(string),
						Upgrade:     entry["upgrade"].(bool),
					}
					if author, authorOk := entry["new_member"].(map[string]interface{}); authorOk {
						event.NewMember = yt_stats.ChatUser{
							AuthorName:       author["author_name"].(string),
							AuthorId:         author["author_id"].(string),
							AuthorChannelUrl: author["author_channel_url"].(string),
							ChatOwner:        author["chat_owner"].(bool),
							Moderator:        author["moderator"].(bool),
							Member:           author["member"].(bool),
							Verified:         author["verified"].(bool),
						}
					}
					outbound.ChatEvents[i] = event
				} else if entry["type"] == "membership_milestone" {
					event := yt_stats.ChatMemberMilestone{
						Id:          entry["id"].(string),
						Type:        "membership_milestone",
						PublishedAt: entry["published_at"].(string),
						Message:     entry["message"].(string),
						UserComment: entry["user_comment"].(string),
						Level:       entry["level"].(string),
						Months:      int(entry["months"].(float64)),
					}
					if author, authorOk := entry["member"].(map[string]interface{}); authorOk {
						event.Member = yt_stats.ChatUser{
							AuthorName:       author["author_name"].(string),
							AuthorId:         author["author_id"].(string),
							AuthorChannelUrl: author["author_channel_url"].(string),
							ChatOwner:        author["chat_owner"].(bool),
							Moderator:        author["moderator"].(bool),
							Member:           author["member"].(bool),
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
							Member:           author["member"].(bool),
							Verified:         author["verified"].(bool),
						}
					}
					outbound.ChatEvents[i] = event
				} else {
					outbound.ChatEvents[i] = nil
				}
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
	expected := fmt.Sprintf(`{"quota_usage":5,"status_code":%d,"status_message":"liveChatEnded"}`,
		http.StatusForbidden)
	if strings.Trim(rr.Body.String(), "\n") != expected {
		t.Errorf("handler returned wrong body: expected %v actually %v", expected, rr.Body.String())
	}
}

func TestChatHandlerUnsupportedType(t *testing.T) {
	unsupportedRequestType(t, yt_stats.ChatHandler, "/ytstats/v1/chat/", "PUT")
}
