/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

func TestNewErrorEvent(t *testing.T) {
	assert := assert.New(t)

	/// stuff
	ee := NewErrorEvent(
		Fatal,
		fmt.Errorf("not a test"),
		OptErrorEventState(&http.Request{Method: "POST"}),
	)
	assert.Equal(Fatal, ee.GetFlag())
	assert.Equal("not a test", ee.Err.Error())
	assert.NotNil(ee.State)
	assert.Equal(ee.Restricted, false)
	assert.Equal("POST", ee.State.(*http.Request).Method)

	buf := new(bytes.Buffer)
	tf := TextOutputFormatter{
		NoColor: true,
	}

	ee.WriteText(tf, buf)
	assert.Equal("not a test", buf.String())

	contents, err := json.Marshal(ee.Decompose())
	assert.Nil(err)
	assert.Contains(string(contents), "not a test")

	ee = NewErrorEvent(Fatal, ex.New("this is only a test"), OptErrorEventRestricted(true))
	contents, err = json.Marshal(ee.Decompose())
	assert.Nil(err)
	assert.Equal(ee.Restricted, true)
	assert.Contains(string(contents), "this is only a test")
}

func TestErrorEventListener(t *testing.T) {
	assert := assert.New(t)

	ee := NewErrorEvent(Fatal, fmt.Errorf("only a test"))

	var didCall bool
	ml := NewErrorEventListener(func(ctx context.Context, e ErrorEvent) {
		didCall = true
	})

	ml(context.Background(), ee)
	assert.True(didCall)
}

func TestScopedErrorEventListener(t *testing.T) {
	testCases := []struct {
		scopes			*Scopes
		enabledContexts		[]context.Context
		disabledContexts	[]context.Context
	}{
		{
			scopes:	NewScopes("*"),
			enabledContexts: []context.Context{
				WithPath(context.Background(), "test0", "test1"),
				WithPath(context.Background(), "test0"),
				WithPath(context.Background(), "test1"),
			},
		},
		{
			scopes:	NewScopes("-*"),
			disabledContexts: []context.Context{
				WithPath(context.Background(), "test0", "test1"),
				WithPath(context.Background(), "test0"),
				WithPath(context.Background(), "test1"),
			},
		},
		{
			scopes:	NewScopes("test0/test1"),
			enabledContexts: []context.Context{
				WithPath(context.Background(), "test0", "test1"),
			},
			disabledContexts: []context.Context{
				WithPath(context.Background(), "test0"),
				WithPath(context.Background(), "test0", "test2"),
				WithPath(context.Background(), "test0", "test1", "test2"),
				WithPath(context.Background(), "test1"),
			},
		},
	}

	for _, testCase := range testCases {
		name := fmt.Sprintf("Scope '%s'", testCase.scopes.String())
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			for _, ctx := range testCase.enabledContexts {
				ee := NewErrorEvent(Fatal, fmt.Errorf("only a test"))

				var didCall bool
				ml := NewScopedErrorEventListener(func(ctx context.Context, e ErrorEvent) {
					didCall = true
				}, testCase.scopes)

				ml(ctx, ee)
				assert.True(didCall, strings.Join(GetPath(ctx), "/"))
			}

			for _, ctx := range testCase.disabledContexts {
				ee := NewErrorEvent(Fatal, fmt.Errorf("only a test"))

				var didCall bool
				ml := NewScopedErrorEventListener(func(ctx context.Context, e ErrorEvent) {
					didCall = true
				}, testCase.scopes)

				ml(ctx, ee)
				assert.False(didCall, strings.Join(GetPath(ctx), "/"))
			}
		})
	}
}
