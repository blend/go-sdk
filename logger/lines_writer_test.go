/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
)

func TestLinesWriterWrite(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		contents	[]byte
		expectedBuf	string
		expectedCount	int
	}{
		{
			contents:	[]byte("hello\nworld\n!"),
			expectedBuf:	"helloworld!",
			expectedCount:	3,
		},
		{
			contents:	[]byte("asdfsafdsaf"),
			expectedBuf:	"asdfsafdsaf",
			expectedCount:	1,
		},
		{
			contents:	[]byte("\n\n\n\n"),
			expectedBuf:	"",
			expectedCount:	5,
		},
		{
			contents:	[]byte("trailing\n"),
			expectedBuf:	"trailing",
			expectedCount:	2,
		},
		{
			contents:	[]byte(""),
			expectedBuf:	"",
			expectedCount:	1,
		},
	}

	for _, test := range tests {
		var buf bytes.Buffer
		ws := &writerSpy{writer: &buf}
		lw := logger.NewLinesWriter(ws)

		n, err := lw.Write(test.contents)
		assert.Nil(err)
		assert.Equal(len(test.contents), n)
		assert.Equal(test.expectedBuf, buf.String())
		assert.Equal(test.expectedCount, ws.count)
	}
}

type writerSpy struct {
	writer	io.Writer
	count	int
}

func (ws *writerSpy) Write(p []byte) (int, error) {
	ws.count++
	return ws.writer.Write(p)
}
