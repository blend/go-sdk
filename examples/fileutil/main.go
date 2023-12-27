/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/blend/go-sdk/fileutil"
)

func main() {
	go fileutil.Watch(context.Background(), "file.yml", func(f *os.File) error {
		defer f.Close()
		fmt.Printf("file changed\n")
		return nil
	})

	select {}
}
