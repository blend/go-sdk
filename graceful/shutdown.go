/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package graceful

// Shutdown racefully stops a set hosted processes based on SIGINT or SIGTERM received from the os.
// It will return any errors returned by Start() that are not caused by shutting down the server.
// A "Graceful" processes *must* block on start.
func Shutdown(hosted ...Graceful) error {
	return ShutdownBySignal(hosted,
		OptDefaultShutdownSignal(),
	)
}
