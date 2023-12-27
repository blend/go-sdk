/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestTimeout(t *testing.T) {
	for _, tc := range []struct {
		Name	string
		Timeout	time.Duration
		Action	Action
		Status	int
	}{
		{
			Name:		"panic",
			Timeout:	time.Minute,
			Action: func(_ *Ctx) Result {
				panic("test")
			},
			Status:	http.StatusInternalServerError,
		},
		{
			Name:		"long action",
			Timeout:	time.Microsecond,
			Action: func(r *Ctx) Result {
				<-r.Context().Done()
				return NoContent
			},
			Status:	http.StatusServiceUnavailable,
		},
		{
			Name:		"short action",
			Timeout:	time.Minute,
			Action: func(_ *Ctx) Result {
				return NoContent
			},
			Status:	http.StatusNoContent,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			assert := assert.New(t)

			app := MustNew(
				OptBindAddr(DefaultMockBindAddr),
				OptUse(WithTimeout(tc.Timeout)),
			)
			app.GET("/endpoint", tc.Action)

			ts := httptest.NewServer(app)
			defer ts.Close()

			res, err := http.Get(ts.URL + "/endpoint")
			assert.Nil(err)
			assert.Nil(res.Body.Close())
			assert.Equal(tc.Status, res.StatusCode)
		})
	}
}
