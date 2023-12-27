/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

import "context"

// Listener is a function that can be triggered by events.
type Listener func(context.Context, Event)
