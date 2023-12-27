/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"context"

	"github.com/blend/go-sdk/webutil"
)

// OptContext sets the request context.
func OptContext(ctx context.Context) Option {
	return RequestOption(webutil.OptContext(ctx))
}
