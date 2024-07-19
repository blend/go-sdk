/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*
Package vault implements a high throughput vault client.

It also provides helpers for reading and writing objects to vault key value stores.

# Mock and Testing Examples

Very often you will need to mock the vault client in your code so you don't reach out to and actual vault instance during tests.
Before writing tests, however, you should make sure that any references to the vault client do so through the `vault.Client` interface, not a concrete type like `*vault.APIClient`.

Then, in your tests, you can create a new mock:

	type clientMock struct {
		vault.Client // embed the vault client interface to satisfy the interface requirements.
	}
	// implement a specific method you need to mock
	func (clientMock) Get(_ context.Context, path string, opts ...vault.CallOption) (vault.Values, error) {
		return vault.Values{ "foo": "bar"}, nil
	}

This will then let you pass `new(clientMock)` to anywhere you need to set a `vault.Client`
*/
package vault // import "github.com/blend/go-sdk/vault"
