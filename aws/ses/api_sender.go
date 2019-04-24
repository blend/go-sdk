package ses

import (
	"context"

	awsutil "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awsSes "github.com/aws/aws-sdk-go/service/ses"

	"github.com/blend/go-sdk/email"
	"github.com/blend/go-sdk/ex"
)

var (
	_ Sender = (*APISender)(nil)
)

// New returns a new sender.
func New(session *session.Session) email.Sender {
	return &APISender{
		Session: session,
		Client:  awsSes.New(session),
	}
}

// APISender is an aws ses email sender.
type APISender struct {
	Session *session.Session
	Client  *awsSes.SES
}

// Send sends a message.
func (s *APISender) Send(ctx context.Context, m email.Message) error {
	if s.Client == nil {
		return nil
	}
	if err := m.Validate(); err != nil {
		return err
	}
	input := &awsSes.SendEmailInput{
		Source: &m.From,
		Destination: &awsSes.Destination{
			ToAddresses:  awsutil.StringSlice(m.To),
			CcAddresses:  awsutil.StringSlice(m.CC),
			BccAddresses: awsutil.StringSlice(m.BCC),
		},
		Message: &awsSes.Message{
			Subject: &awsSes.Content{
				Charset: &defaultCharset,
				Data:    &m.Subject,
			},
			Body: &awsSes.Body{},
		},
	}

	if len(m.HTMLBody) > 0 {
		input.Message.Body.Html = &awsSes.Content{
			Charset: &defaultCharset,
			Data:    &m.HTMLBody,
		}
	}

	if len(m.TextBody) > 0 {
		input.Message.Body.Text = &awsSes.Content{
			Charset: &defaultCharset,
			Data:    &m.TextBody,
		}
	}

	_, err := s.Client.SendEmailWithContext(ctx, input)
	return ex.New(err)
}
