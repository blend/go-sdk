/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
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
