/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package ex

// InnerProvider is a type that returns an inner error.
type InnerProvider interface {
	Inner() error
}
