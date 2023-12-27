/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package ex

// InnerProvider is a type that returns an inner error.
type InnerProvider interface {
	Inner() error
}
