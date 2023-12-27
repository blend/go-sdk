/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*
Package ex provides the foundations for error handling in the SDK tree.

To create an error that includes a given string class and stack trace:

	err := ex.New("this is a structured error")
	...
	fmt.Println(ex.ErrStackTrace(err))

When in doubt, wrap any errors from non-sdk methods with an exception:

	res, err := http.Get(...)
	if err != nil {
		return nil, ex.New(err) // stack trace will originate from this call
	}

To create an error from a known error class, that can be used later to check the type of the error:

	var ErrTooManyFoos ex.Class = "too many foos"
	...
	err := ex.New(ErrTooManyFoos)
	...
	if ex.Is(err, ErrTooManyFoos) { // we can now verify the type of the err with `ex.Is(err, class)`
		fmt.Println("We did too many foos!")
	}

We can pass other options to the `ex.New(...)` constructor, such as setting an inner error:

	err := ex.New(ErrValidation, ex.OptInner(err))
	...
	if ex.Is(err, ErrValidation) {
		fmt.Printf("validation error: %v\n", ex.ErrInner(err))
	}
*/
package ex // import "github.com/blend/go-sdk/ex"
