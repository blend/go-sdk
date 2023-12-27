/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package async

import "context"

// ContextAction is an action that is given a context and returns an error.
type ContextAction func(ctx context.Context) error
