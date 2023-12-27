/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
