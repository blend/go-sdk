/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*
Package assert adds helpers to make writing tests easier.

Example Usage:

	func TestFoo(t *testing.T) {
		// create the assertions wrapper
		assert := assert.New(t)

		assert.True(false) // this will fail the test.
	}
*/
package assert
