/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*
Package r2test provides helpers for writing tests involving calls with sdk/r2.

The most common example is to add a mock response as an option to a default set of options.

Lets say we have a wrapping helper client:

	type APIClient struct {
		RemoteURL string
		...
		Defaults []r2.Option
	}

	func (a APIClient) GetFoos(ctx context.Context) (output []Foo, err error) {
		_, err = r2.New(a.RemoteURL, append(a.Defaults, r2.OptContext(ctx))...).JSON(&output)
		return
	}

During tests we can add a default option:

	mockedResponse := `[{"is":"a foo"}]`
	a := APIClient{ Remote: "http://test.invalid", Defaults: []r2.Option{r2test.OptMockResponseString(mockedResponse)} }
	foos, err := a.GetFoos(context.TODO())
	...

We will now return the mocked response instead of reaching out to the remote for the call.
*/
package r2test // import "github.com/blend/go-sdk/r2/r2test"
