/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

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
