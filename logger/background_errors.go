/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

// BackgroundErrors reads errors from a channel and logs them as errors.
//
// You should call this method with it's own goroutine:
//
//	go logger.BackgroundErrors(log, flushErrors)
func BackgroundErrors(log ErrorReceiver, errors <-chan error) {
	if !IsLoggerSet(log) {
		return
	}
	var err error
	for {
		err = <-errors
		if err != nil {
			log.Error(err)
		}
	}
}
