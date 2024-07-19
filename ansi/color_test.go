/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package ansi

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestColorApply(t *testing.T) {
	assert := assert.New(t)

	escapedBlack := ColorBlack.Normal()
	assert.Equal("\033[0;"+string(ColorBlack), escapedBlack)

	appliedBlack := ColorBlack.Apply("test")
	assert.Equal(ColorBlack.Normal()+"test"+ColorReset, appliedBlack)
}

func TestColors(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(ColorBlack.Apply("foo"), Black("foo"))
	assert.Equal(ColorRed.Apply("foo"), Red("foo"))
	assert.Equal(ColorGreen.Apply("foo"), Green("foo"))
	assert.Equal(ColorYellow.Apply("foo"), Yellow("foo"))
	assert.Equal(ColorBlue.Apply("foo"), Blue("foo"))
	assert.Equal(ColorPurple.Apply("foo"), Purple("foo"))
	assert.Equal(ColorCyan.Apply("foo"), Cyan("foo"))
	assert.Equal(ColorWhite.Apply("foo"), White("foo"))
	assert.Equal(ColorLightBlack.Apply("foo"), LightBlack("foo"))
	assert.Equal(ColorLightRed.Apply("foo"), LightRed("foo"))
	assert.Equal(ColorLightGreen.Apply("foo"), LightGreen("foo"))
	assert.Equal(ColorLightYellow.Apply("foo"), LightYellow("foo"))
	assert.Equal(ColorLightBlue.Apply("foo"), LightBlue("foo"))
	assert.Equal(ColorLightPurple.Apply("foo"), LightPurple("foo"))
	assert.Equal(ColorLightCyan.Apply("foo"), LightCyan("foo"))
	assert.Equal(ColorLightWhite.Apply("foo"), LightWhite("foo"))

	assert.Equal(ColorRed.Bold()+"foo"+ColorReset, Bold(ColorRed, "foo"))
	assert.Equal(ColorRed.Underline()+"foo"+ColorReset, Underline(ColorRed, "foo"))
}
