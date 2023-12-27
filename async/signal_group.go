/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package async

import "sync"

// SignalGroup is a wait group but it awaits on a signal.
type SignalGroup struct {
	wg sync.WaitGroup
}

// Add adds delta.
func (sg *SignalGroup) Add(delta int) {
	sg.wg.Add(delta)
}

// Done decrements delta.
func (sg *SignalGroup) Done() {
	sg.wg.Done()
}

// Wait returns a channel you can select from.
func (sg *SignalGroup) Wait() <-chan struct{} {
	finished := make(chan struct{})
	go func() {
		sg.wg.Wait()
		close(finished)
	}()
	return finished
}
