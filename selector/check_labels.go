/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
