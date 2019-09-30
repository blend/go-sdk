package slack

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/blend/go-sdk/r2"

	"net/http"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/r2"
)

const (
	// ErrNon200 is the exception class when a non-200 is returned from slack.
	ErrNon200 = "slack; non-200 status code returned from remote"
)

var (
	_ Sender = (*WebhookSender)(nil)
)

// New creates a new webhook sender.
func New(cfg Config) *WebhookSender {
	return &WebhookSender{
		Transport: &http.Transport{},
		Config:    cfg,
	}
}

// WebhookSender sends slack webhooks.
type WebhookSender struct {
	Transport       *http.Transport
	RequestDefaults []r2.Option
	Config          Config
}

// Defaults returns default message options.
func (whs WebhookSender) Defaults() []MessageOption {
	return []MessageOption{
		WithUsernameOrDefault(whs.Config.Username),
		WithChannelOrDefault(whs.Config.Channel),
		WithIconEmojiOrDefault(whs.Config.IconEmoji),
		WithIconURLOrDefault(whs.Config.IconURL),
	}
}

// Send sends a slack hook.
func (whs WebhookSender) Send(ctx context.Context, message Message) error {
	messageWithDefaults := ApplyMessageOptions(message, whs.Defaults()...)

	options := append(whs.RequestDefaults, r2.OptTransport(whs.Transport), r2.OptJSONBody(messageWithDefaults))

	res, err := r2.New(whs.Config.Webhook, options...).Do()
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode > http.StatusOK {
		contents, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return ex.New(err)
		}
		return ex.New(ErrNon200, ex.OptMessage(string(contents)))
	}
	return nil
}

// SendAndReadResponse sends a slack hook and returns the deserialized response
func (whs WebhookSender) SendAndReadResponse(ctx context.Context, message Message) (*PostMessageResponse, error) {
	messageWithDefaults := ApplyMessageOptions(message, whs.Defaults()...)

	options := append(whs.RequestDefaults, r2.OptTransport(whs.Transport), r2.OptJSONBody(messageWithDefaults))

	res, err := r2.New(whs.Config.Webhook, options...).Do()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var contents PostMessageResponse
	err = json.NewDecoder(res.Body).Decode(&contents)
	if err != nil {
		return nil, ex.New(err)
	}

	if res.StatusCode > http.StatusOK {
		return &contents, ex.New(ErrNon200, ex.OptMessagef("%#v", contents))
	}

	return &contents, nil
}

// PostMessage posts a basic message to a given channel.
func (whs WebhookSender) PostMessage(channel, messageText string, options ...MessageOption) error {
	message := Message{
		Channel: channel,
		Text:    messageText,
	}
	for _, option := range options {
		option(&message)
	}
	return whs.Send(context.Background(), message)
}

// PostMessageAndReadResponse posts a basic message to a given channel and returns the deserialized response
func (whs WebhookSender) PostMessageAndReadResponse(channel, messageText string, options ...MessageOption) (*PostMessageResponse, error) {
	message := Message{
		Channel: channel,
		Text:    messageText,
	}
	for _, option := range options {
		option(&message)
	}
	return whs.SendAndReadResponse(context.Background(), message)
}

// PostMessageContext posts a basic message to a given chanel with a given context.
func (whs WebhookSender) PostMessageContext(ctx context.Context, channel, messageText string, options ...MessageOption) error {
	message := Message{
		Channel: channel,
		Text:    messageText,
	}
	for _, option := range options {
		option(&message)
	}
	return whs.Send(ctx, message)
}
