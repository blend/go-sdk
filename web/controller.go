/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package web

// Controller is an interface for controller objects.
/*
The primary concern of a controller is to register routes that correspond to the
actions the controller implements.

Routes are registered in order, and cannot collide with eachother.

Controllers should also register any views or additional resources they need
at the time of registration.
*/
type Controller interface {
	Register(app *App)
}
