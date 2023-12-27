/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package graceful

// Graceful is a server that can start and stop.
type Graceful interface {
	Start() error	// this call must block
	Stop() error
}
