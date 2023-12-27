/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

import (
	"github.com/blend/go-sdk/ansi"
)

var (
	// DefaultFlagTextColors is the default color for each known flag.
	DefaultFlagTextColors = map[string]ansi.Color{
		Info:    ansi.ColorLightWhite,
		Debug:   ansi.ColorLightYellow,
		Warning: ansi.ColorLightYellow,
		Error:   ansi.ColorRed,
		Fatal:   ansi.ColorRed,
	}

	// DefaultFlagTextColor is the default flag color.
	DefaultFlagTextColor = ansi.ColorLightWhite
)

// FlagTextColor returns the color for a flag.
func FlagTextColor(flag string) ansi.Color {
	if color, hasColor := DefaultFlagTextColors[flag]; hasColor {
		return color
	}
	return DefaultFlagTextColor
}
