/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package logger

// EventContext is a wrapping context for events.
// It is used when a sub-context triggers or writes an event.
type EventContext struct {
	Event
	ContextPath	[]string
}
