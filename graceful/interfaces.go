/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package graceful

// Graceful is a server that can start and stop.
type Graceful interface {
	Start() error // this call must block
	Stop() error
}
