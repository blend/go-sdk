/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package async

import "context"

// WorkAction is an action handler for a queue.
type WorkAction func(context.Context, interface{}) error
