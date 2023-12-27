/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/email"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/stringutil"
)

type flagStrings []string

func (fs *flagStrings) Values() []string {
	return []string(*fs)
}

func (fs *flagStrings) String() string {
	return "an array of strings"
}

func (fs *flagStrings) Set(value string) error {
	*fs = append(*fs, value)
	return nil
}

func main() {
	var to flagStrings
	var (
		from    = flag.String("from", "noreply@example.org", "The message `from` address")
		subject = flag.String("subject", "", "The message `subject`")
		body    = flag.String(`body`, "", "The message `body`")
	)
	flag.Var(&to, "to", "The message `to` address(es), can be more than one.")
	flag.Parse()

	log := logger.Prod()

	var sender email.SMTPSender
	if _, err := configutil.Read(&sender); !configutil.IsIgnored(err) {
		logger.FatalExit(err)
	}

	message := email.Message{
		From:     *from,
		To:       to.Values(),
		Subject:  *subject,
		TextBody: *body,
	}

	log.Infof("using smtp host:     %s", sender.Host)
	log.Infof("using smtp username: %s", sender.PlainAuth.Username)
	log.Infof("using smtp port:     %s", sender.PortOrDefault())
	log.Infof("using message from address:  %s", *from)
	log.Infof("using message to addresses:  %s", stringutil.CSV(to.Values()))
	log.Infof("using message subject:       %s", *subject)

	if err := sender.Send(context.Background(), message); err != nil {
		logger.FatalExit(err)
	}
	fmt.Println("message sent!")
}
