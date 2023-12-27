/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package logger

import "github.com/blend/go-sdk/ansi"

// TextFormatter is a type that can format text output.
type TextFormatter interface {
	Colorize(string, ansi.Color) string
}
