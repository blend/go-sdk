/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package timeutil

import "time"

// Extra time formats
const (
	TimeFormatCompleteDate24TimeUTC      = "2006-01-02 15:04Z07:00"
	TimeFormatCompleteDate24TimeUTCSlash = "2006/01/02 15:04Z07:00"
	TimeFormatCompleteDate               = "2006-01-02"
	TimeFormatCompleteDateSlash          = "2006/01/02"
	TimeFormatKitchen24                  = "15:04"
)

var (
	// DefaultTimeFormats are the default time formats used in parsing.
	// They are ordered from _most_ specific to least specific.
	DefaultTimeFormats = []string{
		time.RFC3339Nano,
		time.RFC3339,
		TimeFormatCompleteDate24TimeUTC,
		TimeFormatCompleteDate24TimeUTCSlash,
		TimeFormatCompleteDate,
		TimeFormatCompleteDateSlash,
		TimeFormatKitchen24,
		time.Kitchen,
	}
)

// ParseTimeDefaults parses a time string with the default time formats.
func ParseTimeDefaults(timeStr string) (time.Time, error) {
	return ParseTime(timeStr, DefaultTimeFormats...)
}

// ParseTime parses a time string with a given set of formats.
func ParseTime(timeStr string, timeFormats ...string) (output time.Time, err error) {
	for _, timeFormat := range timeFormats {
		output, err = time.Parse(timeFormat, timeStr)
		if err == nil && !output.IsZero() {
			return
		}
	}
	return
}
