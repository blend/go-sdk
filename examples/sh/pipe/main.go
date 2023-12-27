/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package main

import (
	"fmt"
	"os"

	"github.com/blend/go-sdk/sh"
)

func main() {

	if err := sh.Pipe(sh.MustCmds("yes", "head -n 5")...); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	if err := sh.Pipe(sh.MustCmds("cat /dev/urandom", "head -c 32", "base64")...); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

}
