/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package reflectutil

import "reflect"

// Type returns the reflect type following pointers.
func Type(obj interface{}) reflect.Type {
	t := reflect.TypeOf(obj)
	for t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface {
		t = t.Elem()
	}

	return t
}
