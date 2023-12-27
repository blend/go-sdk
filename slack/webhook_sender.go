/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package slack

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/r2"
	"github.com/blend/go-sdk/webutil"
)

const (
	// ErrNon200 is the exception class when a non-200 is returned from slack.
	ErrNon200 ex.Class = "slack; non-200 status code returned from remote"
)

var (
	_ Sender = (*WebhookSender)(nil)
)

// New creates a new webhook sender.
func New(cfg Config) *WebhookSender {
	return &WebhookSender{
		Transport:	new(http.Transport),
		Config:		cfg,
	}
}

// WebhookSender sends slack webhooks.
type WebhookSender struct {
	Transport	*http.Transport
	RequestDefaults	[]r2.Option
	Config		Config
}

// MessageDefaults returns default message options.
func (whs WebhookSender) MessageDefaults() []MessageOption {
	return []MessageOption{
		OptMessageUsernameOrDefault(whs.Config.Username),
		OptMessageChannelOrDefault(whs.Config.Channel),
		OptMessageIconEmojiOrDefault(whs.Config.IconEmoji),
		OptMessageIconURLOrDefault(whs.Config.IconURL),
	}
}

func (whs WebhookSender) send(ctx context.Context, message Message) (*http.Response, error) {
	message = ApplyMessageOptions(message, whs.MessageDefaults()...)
	options := append(whs.RequestDefaults,
		r2.OptPost(),
		r2.OptTransport(whs.Transport),
		r2.OptJSONBody(message),
		r2.OptHeaderValue(webutil.HeaderContentType, "application/json"),
	)
	return r2.New(whs.Config.Webhook, options...).Do()
}

// Send sends a slack hook.
func (whs WebhookSender) Send(ctx context.Context, message Message) error {
	res, err := whs.send(ctx, message)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if statusCode := res.StatusCode; statusCode < http.StatusOK || statusCode > 299 {
		contents, _ := io.ReadAll(res.Body)
		return ex.New(ErrNon200, ex.OptMessage(string(contents)))
	}
	return nil
}

// SendAndReadResponse sends a slack hook and returns the deserialized response
func (whs WebhookSender) SendAndReadResponse(ctx context.Context, message Message) (*PostMessageResponse, error) {
	res, err := whs.send(ctx, message)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var contents PostMessageResponse
	err = json.NewDecoder(res.Body).Decode(&contents)
	if err != nil {
		return nil, ex.New(err)
	}
	if statusCode := res.StatusCode; statusCode < http.StatusOK || statusCode > 299 {
		return &contents, ex.New(ErrNon200, ex.OptMessagef("%#v", contents))
	}
	return &contents, nil
}

// PostMessage posts a basic message to a given channel.
func (whs WebhookSender) PostMessage(channel, messageText string, options ...MessageOption) error {
	message := Message{
		Channel:	channel,
		Text:		messageText,
	}
	for _, option := range options {
		option(&message)
	}
	return whs.Send(context.Background(), message)
}

// PostMessageAndReadResponse posts a basic message to a given channel and returns the deserialized response
func (whs WebhookSender) PostMessageAndReadResponse(channel, messageText string, options ...MessageOption) (*PostMessageResponse, error) {
	message := Message{
		Channel:	channel,
		Text:		messageText,
	}
	for _, option := range options {
		option(&message)
	}
	return whs.SendAndReadResponse(context.Background(), message)
}

// PostMessageContext posts a basic message to a given chanel with a given context.
func (whs WebhookSender) PostMessageContext(ctx context.Context, channel, messageText string, options ...MessageOption) error {
	message := Message{
		Channel:	channel,
		Text:		messageText,
	}
	for _, option := range options {
		option(&message)
	}
	return whs.Send(ctx, message)
}

// PostMessageAndReadResponseContext posts a basic message to a given channel and returns the deserialized response
func (whs WebhookSender) PostMessageAndReadResponseContext(ctx context.Context, channel, messageText string, options ...MessageOption) (*PostMessageResponse, error) {
	message := Message{
		Channel:	channel,
		Text:		messageText,
	}
	for _, option := range options {
		option(&message)
	}
	return whs.SendAndReadResponse(ctx, message)
}
