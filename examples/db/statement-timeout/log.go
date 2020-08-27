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
