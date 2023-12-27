/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

/*
Package logger is our high throughput event bus.

It has two main modes of output; text and json, and allows multiple listeners to be triggerd for a given logger event.

The output is governed by the `LOG_FORMAT` environment variable. Text output is the default, which
is great for reading locally, but is less than optimal for search and automated ingestion. In
production systems, `LOG_FORMAT=json` is recommended.
*/
package logger	// import "github.com/blend/go-sdk/logger"
