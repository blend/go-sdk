/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package sentry

import (
	"context"

	"github.com/blend/go-sdk/logger"
)

// Sender is the type that the sentry client ascribes to.
type Sender interface {
	Notify(context.Context, logger.ErrorEvent)
}
