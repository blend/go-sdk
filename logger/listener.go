/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package logger

import "context"

// Listener is a function that can be triggered by events.
type Listener func(context.Context, Event)
