/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package ex

// ClassProvider is a type that can return an exception class.
type ClassProvider interface {
	Class() error
}
