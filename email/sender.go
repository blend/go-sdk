/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package email

import (
	"context"
)

// Sender is a generalized sender.
type Sender interface {
	Send(context.Context, Message) error
}
