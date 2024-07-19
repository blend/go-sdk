/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package assert

import (
	"encoding/json"
	"fmt"
	"strings"
)

// NewFailure creates a new failure.
func NewFailure(message string, userMessageComponents ...interface{}) Failure {
	return Failure{
		Message:     message,
		UserMessage: fmt.Sprint(userMessageComponents...),
		CallerInfo:  callerInfoStrings(callerInfo()),
	}
}

// Failure is an assertion failure.
type Failure struct {
	Message     string   `json:"message"`
	UserMessage string   `json:"userMessage,omitempty"`
	CallerInfo  []string `json:"callerInfo"`
}

// Error implements error.
func (f Failure) Error() string {
	return f.Message
}

// Text returns the text format of the failure.
func (f Failure) Text() string {
	errorTrace := strings.Join(f.CallerInfo, "\n\t")
	if len(errorTrace) == 0 {
		errorTrace = "Unknown"
	}
	assertionFailedLabel := color("Assertion Failed!", RED)
	locationLabel := color("Assert Location", GRAY)
	assertionLabel := color("Assertion", GRAY)
	messageLabel := color("Message", GRAY)
	if f.UserMessage != "" {
		errorFormat := "%s\n%s\n%s:\n\t%s\n%s:\n\t%s\n%s:\n\t%s\n\n"
		return fmt.Sprintf(errorFormat, "", assertionFailedLabel, locationLabel, errorTrace, assertionLabel, f.Message, messageLabel, f.UserMessage)
	}
	errorFormat := "%s\n%s\n%s:\n\t%s\n%s:\n\t%s\n\n"
	return fmt.Sprintf(errorFormat, "", assertionFailedLabel, locationLabel, errorTrace, assertionLabel, f.Message)
}

// JSON returns the json format of the failure.
func (f Failure) JSON() string {
	contents, _ := json.Marshal(f)
	return string(contents)
}

// TestString returns a plain text representation of the contents of a Failure, suitable for easy comparison within unit tests.
// Only the first line of `Failure.CallerInfo` is included, if present, as an assertion of the entire call stack would be extremely verbose,
// brittle, and offer little value.
func (f Failure) TestString() string {
	// this separator is unlikely to naturally occur in the formatted components, which aides the caller to disambiguate the failure internal state.
	const separator = ";-- "
	var res = f.Message
	if f.UserMessage != "" {
		if res == "" {
			res = f.UserMessage
		} else {
			res = res + separator + f.UserMessage
		}
	}
	if len(f.CallerInfo) >= 1 {
		if res == "" {
			res = f.CallerInfo[0]
		} else {
			res = res + separator + f.CallerInfo[0]
		}
	}
	return res
}
