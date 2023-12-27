/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package async

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_Future_Complete(t *testing.T) {
	its := assert.New(t)

	shouldComplete := make(chan struct{})
	didComplete := make(chan struct{})
	action := func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return nil
		case <-shouldComplete:
			close(didComplete)
			return nil
		}
	}

	f := Await(context.TODO(), action)

	close(shouldComplete)
	<-didComplete

	its.Nil(f.Complete())
	its.Nil(f.Complete(), "repeat calls should just return")
}

func Test_Future_Cancel(t *testing.T) {
	its := assert.New(t)

	didCancel := make(chan struct{})
	action := func(ctx context.Context) error {
		<-ctx.Done()
		close(didCancel)
		return nil
	}

	its.Nil(Await(context.TODO(), action).Cancel())
	<-didCancel
}

func Test_Future_Cancel_repeat(t *testing.T) {
	its := assert.New(t)

	didCancel := make(chan struct{})
	action := func(ctx context.Context) error {
		<-ctx.Done()
		close(didCancel)
		return nil
	}

	f := Await(context.TODO(), action)
	its.Nil(f.Cancel())
	<-didCancel

	its.Equal(ErrCannotCancel, f.Cancel())
}

func Test_Future_Finished(t *testing.T) {
	its := assert.New(t)

	shouldComplete := make(chan struct{})
	action := func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return context.Canceled
		case <-shouldComplete:
			return nil
		}
	}

	f := Await(context.TODO(), action)
	close(shouldComplete)
	its.Nil(<-f.Finished())
}
