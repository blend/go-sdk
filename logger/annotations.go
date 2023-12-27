/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

// CombineAnnotations combines one or many set of annotations.
func CombineAnnotations(annotations ...Annotations) Annotations {
	output := make(Annotations)
	for _, set := range annotations {
		if set == nil || len(set) == 0 {
			continue
		}
		for key, value := range set {
			output[key] = value
		}
	}
	return output
}

// Annotations are a collection of string name value pairs.
type Annotations map[string]interface{}
