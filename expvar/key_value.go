/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package expvar

// KeyValue represents a single entry in a Map.
type KeyValue struct {
	Key	string
	Value	Var
}
