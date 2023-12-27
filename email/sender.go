/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package email

import (
	"context"
)

// Sender is a generalized sender.
type Sender interface {
	Send(context.Context, Message) error
}
