/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package ex

import "encoding/json"

// Class is a string wrapper that implements `error`.
// Use this to implement constant exception causes.
type Class string

// Class implements `error`.
func (c Class) Error() string {
	return string(c)
}

// MarshalJSON implements json.Marshaler.
func (c Class) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(c))
}
