/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

// EventContext is a wrapping context for events.
// It is used when a sub-context triggers or writes an event.
type EventContext struct {
	Event
	ContextPath []string
}
