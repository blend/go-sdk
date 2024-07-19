/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package email

import "context"

// NewMockSender creates a new mock sender.
func NewMockSender() MockSender {
	return MockSender(make(chan Message))
}

// MockSender is a mocked sender.
type MockSender chan Message

// Send sends a mocked message.
func (ms MockSender) Send(ctx context.Context, m Message) error {
	ms <- m
	return nil
}
