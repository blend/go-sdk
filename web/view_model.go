package web

import (
	"github.com/blend/go-sdk/env"
)

// ViewModel is a wrapping viewmodel.
type ViewModel struct {
	Env       env.Vars
	Status    ViewStatus
	Ctx       *Ctx
	ViewModel interface{}
}

// Wrap returns a ViewModel that wraps a new object.
func (vm ViewModel) Wrap(other interface{}) ViewModel {
	return ViewModel{
		Env:       vm.Env,
		Ctx:       vm.Ctx,
		ViewModel: other,
	}
}
