/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package async

import "context"

// WorkerFinalizer is an action handler for a queue.
type WorkerFinalizer func(context.Context, *Worker) error
