/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package email

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

func TestMessageValidate(t *testing.T) {
	assert := assert.New(t)

	assert.True(ex.Is(Message{}.Validate(), ErrMessageFieldUnset))
	assert.True(ex.Is(Message{
		From: "foo@bar.com",
	}.Validate(), ErrMessageFieldUnset))
	assert.True(ex.Is(Message{
		From: "foo\r@bar.com",
	}.Validate(), ErrMessageFieldNewlines))
	assert.True(ex.Is(Message{
		From: "foo\n@bar.com",
	}.Validate(), ErrMessageFieldNewlines))
	assert.True(ex.Is(Message{
		From: "foo\r\n@bar.com",
	}.Validate(), ErrMessageFieldNewlines))
	assert.True(ex.Is(Message{
		From: "foo@bar.com",
		To:   []string{"moo@bar.com", "bad\n@bar.com"},
	}.Validate(), ErrMessageFieldNewlines))
	assert.True(ex.Is(Message{
		From: "foo@bar.com",
		To:   []string{"moo@bar.com"},
		CC:   []string{"bad\n@bar.com"},
	}.Validate(), ErrMessageFieldNewlines))
	assert.True(ex.Is(Message{
		From: "foo@bar.com",
		To:   []string{"moo@bar.com"},
		CC:   []string{"ok@bar.com"},
		BCC:  []string{"bad\n@bar.com"},
	}.Validate(), ErrMessageFieldNewlines))
	assert.True(ex.Is(Message{
		From:    "foo@bar.com",
		To:      []string{"moo@bar.com"},
		Subject: "this is \n bad",
	}.Validate(), ErrMessageFieldNewlines))
	assert.True(ex.Is(Message{
		From:    "foo@bar.com",
		To:      []string{"moo@bar.com"},
		Subject: "this is \r bad",
	}.Validate(), ErrMessageFieldNewlines))
	assert.True(ex.Is(Message{
		From:    "foo@bar.com",
		To:      []string{"moo@bar.com"},
		Subject: "this is \n\r bad",
	}.Validate(), ErrMessageFieldNewlines))
	assert.True(ex.Is(Message{
		From: "foo@bar.com",
		To:   []string{"moo@bar.com"},
	}.Validate(), ErrMessageFieldUnset))

	assert.Nil(Message{
		From:     "foo@bar.com",
		To:       []string{"moo@bar.com"},
		TextBody: "stuff",
	}.Validate())
}

func TestMessageSerializeJSON(t *testing.T) {
	assert := assert.New(t)

	contents, err := json.Marshal(&Message{})
	assert.Nil(err)
	assert.NotEmpty(contents)
}

func TestMessageSerializeYAML(t *testing.T) {
	assert := assert.New(t)

	contents, err := yaml.Marshal(&Message{})
	assert.Nil(err)
	assert.NotEmpty(contents)
}
