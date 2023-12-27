/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package profanity

// Logger are the methods required on the logger.
type Logger interface {
	Printf(string, ...interface{})
	Errorf(string, ...interface{})
}
