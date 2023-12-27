/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package email

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/ex"
)

var (
	_ Sender = (*SMTPSender)(nil)
)

// SMTPSender is a sender for emails over smtp.
type SMTPSender struct {
	LocalName	string		`json:"localname" yaml:"localname"`
	Host		string		`json:"host" yaml:"host" env:"SMTP_HOST"`
	Port		string		`json:"port" yaml:"port" env:"SMTP_PORT"`
	PlainAuth	SMTPPlainAuth	`json:"plainAuth" yaml:"plainAuth"`
}

// Resolve implements configutil.ConfigResolver.
func (s *SMTPSender) Resolve(ctx context.Context) error {
	return configutil.Resolve(ctx,
		func(ictx context.Context) error { return env.GetVars(ictx).ReadInto(s) },
		s.PlainAuth.Resolve,
	)
}

// IsZero returns if the smtp sender is set or not.
func (s SMTPSender) IsZero() bool {
	return s.Host == ""
}

// PortOrDefault returns a property or a default.
func (s SMTPSender) PortOrDefault() string {
	if s.Port != "" {
		return s.Port
	}
	return "465"
}

// LocalNameOrDefault returns a property or a default.
func (s SMTPSender) LocalNameOrDefault() string {
	if s.LocalName != "" {
		return s.LocalName
	}
	return s.Host
}

// Send sends an email via. smtp.
func (s SMTPSender) Send(ctx context.Context, message Message) error {
	if s.Host == "" {
		return ex.New("smtp host unset")
	}
	if err := message.Validate(); err != nil {
		return err
	}

	tlsConfig := &tls.Config{ServerName: s.Host, InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", s.Host, s.PortOrDefault()), tlsConfig)
	if err != nil {
		return ex.New(err)
	}

	client, err := smtp.NewClient(conn, s.Host)
	if err != nil {
		return ex.New(err)
	}
	defer client.Close()

	if err := client.Hello(s.LocalNameOrDefault()); err != nil {
		return ex.New(err)
	}
	if !s.PlainAuth.IsZero() {
		if err := client.Auth(smtp.PlainAuth(s.PlainAuth.Identity, s.PlainAuth.Username, s.PlainAuth.Password, s.Host)); err != nil {
			return ex.New(err)
		}
	}
	if err := client.Mail(message.From); err != nil {
		return ex.New(err)
	}

	for _, to := range message.To {
		if err := client.Rcpt(to); err != nil {
			return ex.New(err)
		}
	}
	for _, cc := range message.CC {
		if err := client.Rcpt(cc); err != nil {
			return ex.New(err)
		}
	}
	for _, bcc := range message.BCC {
		if err := client.Rcpt(bcc); err != nil {
			return ex.New(err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return ex.New(err)
	}

	// msg data
	bufWriter := bufio.NewWriter(w)
	if _, err := bufWriter.WriteString(fmt.Sprintf("From: %s\r\n", message.From)); err != nil {
		return ex.New(err)
	}
	if err = bufWriter.Flush(); err != nil {
		return ex.New(err)
	}
	for _, to := range message.To {
		if _, err := bufWriter.WriteString(fmt.Sprintf("To: %s\r\n", to)); err != nil {
			return ex.New(err)
		}
	}
	if err = bufWriter.Flush(); err != nil {
		return ex.New(err)
	}
	for _, cc := range message.CC {
		if _, err := bufWriter.WriteString(fmt.Sprintf("Cc: %s\r\n", cc)); err != nil {
			return ex.New(err)
		}
	}
	if err = bufWriter.Flush(); err != nil {
		return ex.New(err)
	}
	if message.Subject != "" {
		if _, err := bufWriter.WriteString("Subject: " + message.Subject + "\r\n"); err != nil {
			return ex.New(err)
		}
	}
	if err = bufWriter.Flush(); err != nil {
		return ex.New(err)
	}

	if message.HTMLBody != "" {
		if _, err := bufWriter.WriteString("MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n"); err != nil {
			return ex.New(err)
		}
		if _, err := bufWriter.WriteString(message.HTMLBody); err != nil {
			return ex.New(err)
		}
	} else if message.TextBody != "" {
		if _, err := bufWriter.WriteString(message.TextBody); err != nil {
			return ex.New(err)
		}
	}
	if err = bufWriter.Flush(); err != nil {
		return ex.New(err)
	}
	if err := w.Close(); err != nil {
		return ex.New(err)
	}

	return ex.New(client.Quit())
}

// SMTPPlainAuth is a auth set for smtp.
type SMTPPlainAuth struct {
	Identity	string	`json:"identity" yaml:"identity"`
	Username	string	`json:"username" yaml:"username" env:"SMTP_USERNAME"`
	Password	string	`json:"password" yaml:"password" env:"SMTP_PASSWORD"`
}

// Resolve implements configutil.ConfigResolver.
func (spa SMTPPlainAuth) Resolve(ctx context.Context) error {
	return env.GetVars(ctx).ReadInto(&spa)	//note(wc); i'm not sure this will always work
}

// IsZero returns if the plain auth is unset.
func (spa SMTPPlainAuth) IsZero() bool {
	return spa.Username == "" && spa.Password == ""
}
