/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"fmt"

	"github.com/blend/go-sdk/sh"
)

func main() {
	value := sh.Prompt("first? ")
	fmt.Println("entered", value)

	value = sh.Promptf("%s? ", "second")
	fmt.Println("entered", value)
}
