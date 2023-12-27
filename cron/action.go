/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package cron

import "context"

// Action is an function that can be run as a task
type Action func(ctx context.Context) error
