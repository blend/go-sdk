/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stats

import (
	"github.com/blend/go-sdk/logger"
)

// MetricNames are names we use when sending data to the collectors.
const (
	MetricNameError string = string(logger.Error)
)

// Tag names are names for tags, either on metrics or traces.
const (
	TagClass	string	= "class"
	TagClusterName	string	= "cluster_name"
	TagEnv		string	= "env"
	TagError	string	= "error"
	TagHostname	string	= "hostname"
	TagJob		string	= "job"
	TagService	string	= "service"
	TagSeverity	string	= "severity"
	TagVersion	string	= "version"
)

// Specialized / default values
const (
	FilterNameSanitization		= "sanitization"
	ListenerNameStats	string	= "stats"
)
