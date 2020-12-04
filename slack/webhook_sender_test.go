package slack

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestWebhookSender(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	config := Config{
		Webhook: ts.URL,
	}

	sender := New(config)
	err := sender.Send(context.TODO(), Message{
		Text: "this is only a test",
	})
	assert.Nil(err)
}

func TestWebhookSenderStatusCode(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	config := Config{
		Webhook: ts.URL,
	}

	sender := New(config)
	assert.NotNil(sender.Send(context.TODO(), Message{
		Text: "this is only a test",
	}))
}

func TestWebhookSenderDefaults(t *testing.T) {
	assert := assert.New(t)

	config := Config{
		Webhook:  "http://foo.com",
		Channel:  "#bot-test",
		Username: "default-test",
	}

	sender := New(config)

	message := Message{
		Text: "this is only a test",
	}

	defaults := sender.MessageDefaults()

	for _, option := range defaults {
		option(&message)
	}

	assert.Equal("this is only a test", message.Text)
	assert.Equal("#bot-test", message.Channel)
	assert.Equal("default-test", message.Username)
}

func TestWebhookSendAndReadResponse(t *testing.T) {
	assert := assert.New(t)

	mockResponse := PostMessageResponse{
		OK:        true,
		Channel:   "#bot-test",
		Timestamp: "1503435956.000247",
		Message: Message{
			Text:     "Here's a message for you",
			Username: "ecto1",
			BotID:    "B19LU7CSY",
			Attachments: []MessageAttachment{
				{
					Text: "This is an attachment",
				},
			},
			Type:      "message",
			SubType:   "bot_message",
			Timestamp: "1503435956.000247",
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(mockResponse)
	}))
	defer ts.Close()

	config := Config{
		Webhook: ts.URL,
	}
	sender := New(config)

	// Test: Successful send should return the response body
	response, err := sender.SendAndReadResponse(context.TODO(), Message{
		Text: "this is only a test",
	})
	assert.Nil(err)
	assert.Equal(mockResponse, *response)
}

func TestWebhookSendAndReadResponseStatusCode(t *testing.T) {
	assert := assert.New(t)

	mockResponse := PostMessageResponse{
		OK:    false,
		Error: "too_many_attachments",
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(mockResponse)
	}))
	defer ts.Close()

	config := Config{
		Webhook: ts.URL,
	}
	sender := New(config)

	// Test: Non-200 http response should cause an error to be returned along with the response body
	response, err := sender.SendAndReadResponse(context.TODO(), Message{
		Text: "this is only a test",
	})
	assert.NotNil(err)
	assert.Equal(mockResponse, *response)
}

func TestPostMessageAndReadResponse(t *testing.T) {
	assert := assert.New(t)

	mockResponse := PostMessageResponse{
		OK:        true,
		Channel:   "#bot-test",
		Timestamp: "1503435956.000247",
		Message: Message{
			Text:     "Here's a message for you",
			Username: "ecto1",
			BotID:    "B19LU7CSY",
			Attachments: []MessageAttachment{
				{
					Text: "This is an attachment",
				},
			},
			Type:      "message",
			SubType:   "bot_message",
			Timestamp: "1503435956.000247",
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var received Message
		err := json.NewDecoder(r.Body).Decode(&received)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			mockResponse.Channel = received.Channel
			mockResponse.Message.Text = received.Text
		}
		_ = json.NewEncoder(w).Encode(mockResponse)
	}))
	defer ts.Close()

	config := Config{
		Webhook: ts.URL,
	}
	sender := New(config)

	// Test: Channel and text parameters should be passed along in the request
	expectedChannel, expectedText := "#test-channel", "Test test"
	response, err := sender.PostMessageAndReadResponse(expectedChannel, expectedText)
	assert.Nil(err)
	assert.Equal(true, response.OK)
	assert.Equal(expectedChannel, response.Channel)
	assert.Equal(expectedText, response.Message.Text)
}
