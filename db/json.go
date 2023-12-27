/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package db

import "encoding/json"

// JSON returns the json representation of a given object for inserts / updates.
func JSON(obj interface{}) interface{} {
	jsonBytes, _ := json.Marshal(obj)
	if result := string(jsonBytes); result != "null" { // explicitly bad.
		return result
	}
	return nil
}
