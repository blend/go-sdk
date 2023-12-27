/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package logger

import (
	"context"
)

// Filter mutates an event.
//
// It should return the modified event, and a bool indicating if we should
// drop the event or not. False means we should continue to log the event
// `true` would indicate we should *not* trigger listeners or write output
// for the given event.
type Filter func(context.Context, Event) (e Event, filter bool)
