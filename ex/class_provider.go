/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package ex

// ClassProvider is a type that can return an exception class.
type ClassProvider interface {
	Class() error
}
