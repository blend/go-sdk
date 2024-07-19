/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

// CombineLabels combines one or many set of fields.
func CombineLabels(labels ...Labels) Labels {
	output := make(Labels)
	for _, set := range labels {
		if set == nil || len(set) == 0 {
			continue
		}
		for key, value := range set {
			output[key] = value
		}
	}
	return output
}

// Labels are a collection of string name value pairs.
type Labels map[string]string
