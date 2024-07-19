/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package cron

import "context"

// Action is an function that can be run as a task
type Action func(ctx context.Context) error
