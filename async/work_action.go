/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package async

import "context"

// WorkAction is an action handler for a queue.
type WorkAction func(context.Context, interface{}) error
