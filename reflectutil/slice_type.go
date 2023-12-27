/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package reflectutil

import "reflect"

// SliceType returns the inner type of a slice following pointers.
func SliceType(collection interface{}) reflect.Type {
	v := reflect.ValueOf(collection)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Len() == 0 {
		t := v.Type()
		for t.Kind() == reflect.Ptr || t.Kind() == reflect.Slice {
			t = t.Elem()
		}
		return t
	}
	v = v.Index(0)
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	return v.Type()
}
