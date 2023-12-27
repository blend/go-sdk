/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package main

import (
	"fmt"

	"github.com/blend/go-sdk/template"
)

func main() {
	t := template.New().WithBody("hello {{ .Var \"foo\"}}").WithVar("foo", "world")
	fmt.Println(t.MustProcessString())
}
