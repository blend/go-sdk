/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"os"

	"github.com/blend/go-sdk/ansi"
)

func main() {

	ansi.Table(os.Stdout,
		[]string{"id", "status", "url"},
		[][]string{
			{"0", "200", "http://google.com"},
			{"1", "200", "http://go.blend.com/foo"},
			{"2", "404", "http://go.blend.com/bar"},
		},
	)

	ansi.TableForSlice(os.Stdout,
		[]struct {
			ID     int
			Status int
			URL    string
		}{
			{0, 200, "http://google.com"},
			{1, 200, "http://go.blend.com/foo"},
			{2, 404, "http://go.blend.com/bar"},
		},
	)
}
