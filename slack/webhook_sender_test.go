/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

func Test_WebhookSender_OK(t *testing.T) {
	its := assert.New(t)

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
	its.Nil(err)
}

func Test_WebhookSender_BadStatusCode_ErrorMessage(t *testing.T) {
	its := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "remote message here\n")
	}))
	defer ts.Close()

	config := Config{
		Webhook: ts.URL,
	}

	sender := New(config)
	err := sender.Send(context.TODO(), Message{
		Text: "this is only a test",
	})
	its.NotNil(err)
	its.True(ex.Is(err, ErrNon200))
	its.Equal("remote message here\n", ex.ErrMessage(err))
}

func Test_WebhookSender_Defaults(t *testing.T) {
	its := assert.New(t)

	config := Config{
		Webhook:	"http://foo.com",
		Channel:	"#bot-test",
		Username:	"default-test",
	}

	sender := New(config)

	message := Message{
		Text: "this is only a test",
	}

	defaults := sender.MessageDefaults()

	for _, option := range defaults {
		option(&message)
	}

	its.Equal("this is only a test", message.Text)
	its.Equal("#bot-test", message.Channel)
	its.Equal("default-test", message.Username)
}

func Test_WebhookSend_DecodeResponse(t *testing.T) {
	its := assert.New(t)

	mockResponse := PostMessageResponse{
		OK:		true,
		Channel:	"#bot-test",
		Timestamp:	"1503435956.000247",
		Message: Message{
			Text:		"Here's a message for you",
			Username:	"ecto1",
			BotID:		"B19LU7CSY",
			Attachments: []MessageAttachment{
				{
					Text: "This is an attachment",
				},
			},
			Type:		"message",
			SubType:	"bot_message",
			Timestamp:	"1503435956.000247",
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
	its.Nil(err)
	its.Equal(mockResponse, *response)
}

func Test_WebhookSender_ReadResponse_BadStatusCode(t *testing.T) {
	its := assert.New(t)

	mockResponse := PostMessageResponse{
		OK:	false,
		Error:	"too_many_attachments",
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
	its.NotNil(err)
	its.Equal(mockResponse, *response)
}

func Test_WebhookSender_PostMessageAndReadResponse(t *testing.T) {
	its := assert.New(t)

	mockResponse := PostMessageResponse{
		OK:		true,
		Channel:	"#bot-test",
		Timestamp:	"1503435956.000247",
		Message: Message{
			Text:		"Here's a message for you",
			Username:	"ecto1",
			BotID:		"B19LU7CSY",
			Attachments: []MessageAttachment{
				{
					Text: "This is an attachment",
				},
			},
			Type:		"message",
			SubType:	"bot_message",
			Timestamp:	"1503435956.000247",
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
	its.Nil(err)
	its.Equal(true, response.OK)
	its.Equal(expectedChannel, response.Channel)
	its.Equal(expectedText, response.Message.Text)
}
