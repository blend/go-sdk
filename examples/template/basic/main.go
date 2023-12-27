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
