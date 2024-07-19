/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package reflectutil

import "strings"

// IsExported returns if a field is exported given its name and capitalization.
func IsExported(fieldName string) bool {
	return fieldName != "" && strings.ToUpper(fieldName)[0] == fieldName[0]
}
