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

	"github.com/blend/go-sdk/sh"
)

func main() {
	value := sh.Prompt("first? ")
	fmt.Println("entered", value)

	value = sh.Promptf("%s? ", "second")
	fmt.Println("entered", value)
}
