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
	value, err := sh.Password("first? ")
	if err != nil {
		sh.Fatal(err)
	}
	fmt.Println("entered", value)

	value, err = sh.Passwordf("%s? ", "second")
	if err != nil {
		sh.Fatal(err)
	}
	fmt.Println("entered", value)
}
