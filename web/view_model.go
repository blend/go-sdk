/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import (
	"github.com/blend/go-sdk/env"
)

// ViewModel is a wrapping viewmodel.
type ViewModel struct {
	Env		env.Vars
	Status		ViewStatus
	Ctx		*Ctx
	ViewModel	interface{}
}

// Wrap returns a ViewModel that wraps a new object.
func (vm ViewModel) Wrap(other interface{}) ViewModel {
	return ViewModel{
		Env:		vm.Env,
		Ctx:		vm.Ctx,
		ViewModel:	other,
	}
}

// State returns a state value.
func (vm ViewModel) State(key string) interface{} {
	if vm.Ctx == nil || vm.Ctx.State == nil {
		return nil
	}
	return vm.Ctx.State.Get(key)
}
