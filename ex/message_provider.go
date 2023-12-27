/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package ex

// MessageProvider is a type that returns a message
type MessageProvider interface {
	Message() string
}
