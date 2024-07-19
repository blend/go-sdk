/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package async

import "context"

// ContextAction is an action that is given a context and returns an error.
type ContextAction func(ctx context.Context) error
