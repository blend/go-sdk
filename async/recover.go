/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package async

import (
	"sync"

	"github.com/blend/go-sdk/ex"
)

/*
Recover runs an action and passes any errors to the given errors channel.

This call blocks, if you need it to be backgrounded, you should call it like:

	go Recover(action, errors)
	<-errors
*/
func Recover(action func() error, errors chan error) {
	defer func() {
		if r := recover(); r != nil && errors != nil {
			errors <- ex.New(r)
		}
	}()

	if err := action(); err != nil && errors != nil {
		errors <- err
	}
}

// RecoverGroup runs a recovery against a specific wait group with an error collector.
// It calls Recover internally.
func RecoverGroup(action func() error, errors chan error, wg *sync.WaitGroup) {
	Recover(func() error {
		if wg != nil {
			defer wg.Done()
		}
		return action()
	}, errors)
}
