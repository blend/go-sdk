/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

import "context"

// EventWithContext is an event with the context it was triggered with.
type EventWithContext struct {
	context.Context
	Event
}
