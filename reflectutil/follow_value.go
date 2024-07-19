/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package reflectutil

import "reflect"

// FollowValue derefs a reflect.Value until it isn't a pointer, but will preseve it's nilness.
func FollowValue(v reflect.Value) interface{} {
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return nil
	}

	val := v
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return val.Interface()
}
