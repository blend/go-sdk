/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package web

// Action is the function signature for controller actions.
type Action func(*Ctx) Result

// PanicAction is a receiver for app.PanicHandler.
type PanicAction func(*Ctx, interface{}) Result
