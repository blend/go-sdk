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
