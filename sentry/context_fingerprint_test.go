/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package sentry

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestContextFingerprint(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(GetFingerprint(context.TODO()))
	assert.Nil(GetFingerprint(context.Background()))
	assert.Nil(GetFingerprint(context.WithValue(context.Background(), contextFingerprintKey{}, 1234)))

	assert.Equal([]string{"foo", "bar"}, GetFingerprint(WithFingerprint(context.Background(), "foo", "bar")))
}
