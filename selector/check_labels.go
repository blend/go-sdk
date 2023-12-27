/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package selector

import "github.com/blend/go-sdk/ex"

// CheckLabels validates all the keys and values for the label set.
func CheckLabels(labels Labels) (err error) {
	for key, value := range labels {
		err = CheckKey(key)
		if err != nil {
			err = ex.New(err, ex.OptMessagef("key: %s", key))
			return
		}
		err = CheckValue(value)
		if err != nil {
			err = ex.New(err, ex.OptMessagef("key: %s, value: %s", key, value))
			return
		}
	}
	return
}
