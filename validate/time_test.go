/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package validate

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

func tv(t time.Time) *time.Time	{ return &t }

func TestTimeBefore(t *testing.T) {
	assert := assert.New(t)

	ts := time.Date(2019, 07, 18, 17, 24, 0, 0, time.UTC)

	var verr error
	verr = Time(tv(ts.Add(-time.Hour))).Before(ts)()
	assert.Nil(verr)
	verr = Time(tv(ts.Add(time.Hour))).Before(ts)()
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal(ErrTimeBefore, ErrCause(verr))
}

func TestTimeBeforeNowUTC(t *testing.T) {
	assert := assert.New(t)

	ts := time.Now().UTC()

	var verr error
	verr = Time(tv(ts.Add(-time.Hour))).BeforeNowUTC()()
	assert.Nil(verr)
	verr = Time(tv(ts.Add(time.Hour))).BeforeNowUTC()()
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal(ErrTimeBefore, ErrCause(verr))
}

func TestTimeAfter(t *testing.T) {
	assert := assert.New(t)

	ts := time.Date(2019, 07, 18, 17, 24, 0, 0, time.UTC)

	var verr error
	verr = Time(tv(ts.Add(time.Hour))).After(ts)()
	assert.Nil(verr)
	verr = Time(tv(ts.Add(-time.Hour))).After(ts)()
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal(ErrTimeAfter, ErrCause(verr))
}

func TestTimeAfterNowUTC(t *testing.T) {
	assert := assert.New(t)

	ts := time.Now().UTC()

	var verr error
	verr = Time(tv(ts.Add(time.Hour))).AfterNowUTC()()
	assert.Nil(verr)
	verr = Time(tv(ts.Add(-time.Hour))).AfterNowUTC()()
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal(ErrTimeAfter, ErrCause(verr))
}

func TestTimeBetween(t *testing.T) {
	assert := assert.New(t)

	a := time.Date(2019, 07, 18, 0, 0, 0, 0, time.UTC)
	b := time.Date(2019, 07, 19, 0, 0, 0, 0, time.UTC)
	c := time.Date(2019, 07, 20, 0, 0, 0, 0, time.UTC)

	var verr error
	verr = Time(&b).Between(a, c)()
	assert.Nil(verr)

	verr = Time(&c).Between(a, b)()
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal(ErrTimeBefore, ErrCause(verr))

	verr = Time(&a).Between(b, c)()
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal(ErrTimeAfter, ErrCause(verr))
}
