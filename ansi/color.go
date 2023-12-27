/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package ansi

// Black applies a given color.
func Black(text string) string {
	return Apply(ColorBlack, text)
}

// Red applies a given color.
func Red(text string) string {
	return Apply(ColorRed, text)
}

// Green applies a given color.
func Green(text string) string {
	return Apply(ColorGreen, text)
}

// Yellow applies a given color.
func Yellow(text string) string {
	return Apply(ColorYellow, text)
}

// Blue applies a given color.
func Blue(text string) string {
	return Apply(ColorBlue, text)
}

// Purple applies a given color.
func Purple(text string) string {
	return Apply(ColorPurple, text)
}

// Cyan applies a given color.
func Cyan(text string) string {
	return Apply(ColorCyan, text)
}

// White applies a given color.
func White(text string) string {
	return Apply(ColorWhite, text)
}

// LightBlack applies a given color.
func LightBlack(text string) string {
	return Apply(ColorLightBlack, text)
}

// LightRed applies a given color.
func LightRed(text string) string {
	return Apply(ColorLightRed, text)
}

// LightGreen applies a given color.
func LightGreen(text string) string {
	return Apply(ColorLightGreen, text)
}

// LightYellow applies a given color.
func LightYellow(text string) string {
	return Apply(ColorLightYellow, text)
}

// LightBlue applies a given color.
func LightBlue(text string) string {
	return Apply(ColorLightBlue, text)
}

// LightPurple applies a given color.
func LightPurple(text string) string {
	return Apply(ColorLightPurple, text)
}

// LightCyan applies a given color.
func LightCyan(text string) string {
	return Apply(ColorLightCyan, text)
}

// LightWhite applies a given color.
func LightWhite(text string) string {
	return Apply(ColorLightWhite, text)
}

// Apply applies a given color.
func Apply(colorCode Color, text string) string {
	return colorCode.Normal() + text + ColorReset
}

// Bold applies a given color as bold.
func Bold(colorCode Color, text string) string {
	return colorCode.Bold() + text + ColorReset
}

// Underline applies a given color as underline.
func Underline(colorCode Color, text string) string {
	return colorCode.Underline() + text + ColorReset
}

// Color represents an ansi color code fragment.
type Color string

// Normal escapes the color for use in the terminal.
func (c Color) Normal() string {
	return "\033[0;" + string(c)
}

// Bold escapes the color for use in the terminal as a bold color.
func (c Color) Bold() string {
	return "\033[1;" + string(c)
}

// Underline escapes the color for use in the terminal as underlined color.
func (c Color) Underline() string {
	return "\033[4;" + string(c)
}

// Apply applies a color to a given string.
func (c Color) Apply(text string) string {
	return Apply(c, text)
}

// Utility Color Codes
const (
	ColorReset = "\033[0m"
)

// Color codes
const (
	ColorDefault Color = "39m"
	ColorBlack   Color = "30m"
	ColorRed     Color = "31m"
	ColorGreen   Color = "32m"
	ColorYellow  Color = "33m"
	ColorBlue    Color = "34m"
	ColorPurple  Color = "35m"
	ColorCyan    Color = "36m"
	ColorWhite   Color = "37m"
)

// BrightColorCodes
const (
	ColorLightBlack  Color = "90m"
	ColorLightRed    Color = "91m"
	ColorLightGreen  Color = "92m"
	ColorLightYellow Color = "93m"
	ColorLightBlue   Color = "94m"
	ColorLightPurple Color = "95m"
	ColorLightCyan   Color = "96m"
	ColorLightWhite  Color = "97m"
)

// BackgroundColorCodes
const (
	ColorBackgroundBlack  Color = "40m"
	ColorBackgroundRed    Color = "41m"
	ColorBackgroundGreen  Color = "42m"
	ColorBackgroundYellow Color = "43m"
	ColorBackgroundBlue   Color = "44m"
	ColorBackgroundPurple Color = "45m"
	ColorBackgroundCyan   Color = "46m"
	ColorBackgroundWhite  Color = "47m"
)

// BackgroundColorCodes
const (
	ColorBackgroundBrightBlack  Color = "100m"
	ColorBackgroundBrightRed    Color = "101m"
	ColorBackgroundBrightGreen  Color = "102m"
	ColorBackgroundBrightYellow Color = "103m"
	ColorBackgroundBrightBlue   Color = "104m"
	ColorBackgroundBrightPurple Color = "105m"
	ColorBackgroundBrightCyan   Color = "106m"
	ColorBackgroundBrightWhite  Color = "107m"
)
