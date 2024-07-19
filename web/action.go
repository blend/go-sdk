/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

// Action is the function signature for controller actions.
type Action func(*Ctx) Result

// PanicAction is a receiver for app.PanicHandler.
type PanicAction func(*Ctx, interface{}) Result
