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
