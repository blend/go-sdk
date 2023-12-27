/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package assert

import "os"

const (
	// RED is the ansi escape code fragment for red.
	RED	= "31"
	// BLUE is the ansi escape code fragment for blue.
	BLUE	= "94"
	// GREEN is the ansi escape code fragment for green.
	GREEN	= "32"
	// YELLOW is the ansi escape code fragment for yellow.
	YELLOW	= "33"
	// WHITE is the ansi escape code fragment for white.
	WHITE	= "37"
	// GRAY is the ansi escape code fragment for gray.
	GRAY	= "90"
)

// OutputFormatFromEnv gets the output format from the env or the default.
func OutputFormatFromEnv() OutputFormat {
	outputFormat := OutputFormatText
	if envOutputFormat := os.Getenv("TEST_OUTPUT_FORMAT"); envOutputFormat != "" {
		outputFormat = OutputFormat(envOutputFormat)
	}
	return outputFormat
}

// OutputFormat is an assertion error output format.
type OutputFormat string

// OutputFormats
const (
	OutputFormatDefault	OutputFormat	= ""
	OutputFormatText	OutputFormat	= "text"
	OutputFormatJSON	OutputFormat	= "json"
)
