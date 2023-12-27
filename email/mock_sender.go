/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
