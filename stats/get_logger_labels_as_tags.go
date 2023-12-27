/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stats

import (
	"context"

	"github.com/blend/go-sdk/logger"
)

// GetLoggerLabelsAsTags reads the logger labels map off the context and
// returns the keys and values formatted as a slice of stats tags.
func GetLoggerLabelsAsTags(ctx context.Context) (tags []string) {
	if labels := logger.GetLabels(ctx); len(labels) > 0 {
		for key, value := range labels {
			tags = append(tags, Tag(key, value))
		}
	}
	return
}
