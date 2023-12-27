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
	"time"
)

type logWriter struct {
	Start time.Time
}

func (lw logWriter) Write(bytes []byte) (int, error) {
	d := time.Since(lw.Start)
	s := float64(d) / float64(time.Second)
	return fmt.Printf("%f %s", s, string(bytes))
}

func newLogWriter() logWriter {
	return logWriter{Start: time.Now().UTC()}
}
