/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package validate

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/uuid"
)

func TestUUIDRequired(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = UUID(nil).Required()()
	assert.NotNil(verr)
	assert.Equal(ErrUUIDRequired, ErrCause(verr))

	var empty uuid.UUID
	verr = UUID(&empty).Required()()
	assert.NotNil(verr)
	assert.Equal(ErrUUIDRequired, ErrCause(verr))

	set := uuid.V4()
	verr = UUID(&set).Required()()
	assert.Nil(verr)
}

func TestUUIDForbidden(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = UUID(nil).Forbidden()()
	assert.Nil(verr)

	var empty uuid.UUID
	verr = UUID(&empty).Forbidden()()
	assert.Nil(verr)

	set := uuid.V4()
	verr = UUID(&set).Forbidden()()
	assert.NotNil(verr)
	assert.Equal(ErrUUIDForbidden, ErrCause(verr))
}

func TestUUIDIsV4(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = UUID(nil).IsV4()()
	assert.NotNil(verr)
	assert.Equal(ErrUUIDV4, ErrCause(verr))

	var empty uuid.UUID
	verr = UUID(&empty).IsV4()()
	assert.NotNil(verr)
	assert.Equal(ErrUUIDV4, ErrCause(verr))

	set := uuid.V4()
	verr = UUID(&set).IsV4()()
	assert.Nil(verr)
}

func TestUUIDIsVersion(t *testing.T) {
	assert := assert.New(t)
	version4 := uuid.V4().Version()

	var verr error
	verr = UUID(nil).IsVersion(version4)()
	assert.NotNil(verr)
	assert.Equal(ErrUUIDVersion, ErrCause(verr))

	var empty uuid.UUID
	verr = UUID(&empty).IsVersion(version4)()
	assert.NotNil(verr)
	assert.Equal(ErrUUIDVersion, ErrCause(verr))

	set := uuid.V4()
	verr = UUID(&set).IsVersion(version4)()
	assert.Nil(verr)
}
