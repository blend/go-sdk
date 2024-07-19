/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package vault

// Client is the general interface for a Secrets client
type Client interface {
	KVClient
	TransitClient
}
