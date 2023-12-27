/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package vault

// Client is the general interface for a Secrets client
type Client interface {
	KVClient
	TransitClient
}
