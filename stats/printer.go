/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stats

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/blend/go-sdk/logger"
)

var (
	_ Collector = (*Printer)(nil)
)

// NewPrinter creates a printer from a given logger.
func NewPrinter(log logger.Log) *Printer {
	p := new(Printer)
	if log != nil {
		p.Output = logger.ShimWriter{
			Context: context.Background(),
			Log:     log,
			EventProvider: func(contents []byte) logger.Event {
				return logger.NewMessageEvent("statsd", strings.TrimSpace(string(contents)))
			},
		}
	} else {
		p.Output = os.Stdout
	}
	return p
}

// Printer is a collector that prints calls to a given writer.
type Printer struct {
	Field struct {
		Namespace   string
		DefaultTags []string
	}
	Output io.Writer
}

// AddDefaultTag adds a default tag.
func (p *Printer) AddDefaultTag(name, value string) {
	p.Field.DefaultTags = append(p.Field.DefaultTags, Tag(name, value))
}

// AddDefaultTags adds default tags.
func (p *Printer) AddDefaultTags(tags ...string) {
	p.Field.DefaultTags = append(p.Field.DefaultTags, tags...)
}

// DefaultTags returns the default tags set.
func (p Printer) DefaultTags() []string {
	return p.Field.DefaultTags
}

// Count implemenents stats.Collector.
func (p Printer) Count(name string, value int64, tags ...string) error {
	return p.writeln("count", name, fmt.Sprint(value), tags...)
}

// Increment implemenents stats.Collector.
func (p Printer) Increment(name string, tags ...string) error {
	return p.writeln("increment", name, "1", tags...)
}

// Gauge implemenents stats.Collector.
func (p Printer) Gauge(name string, value float64, tags ...string) error {
	return p.writeln("gauge", name, fmt.Sprint(value), tags...)
}

// Histogram implemenents stats.Collector.
func (p Printer) Histogram(name string, value float64, tags ...string) error {
	return p.writeln("histogram", name, fmt.Sprint(value), tags...)
}

// Distribution implemenents stats.Collector.
func (p Printer) Distribution(name string, value float64, tags ...string) error {
	return p.writeln("distribution", name, fmt.Sprint(value), tags...)
}

// TimeInMilliseconds implemenents stats.Collector.
func (p Printer) TimeInMilliseconds(name string, value time.Duration, tags ...string) error {
	return p.writeln("timeInMilliseconds", name, value.String(), tags...)
}

func (p Printer) writeln(metricType, name, value string, tags ...string) (err error) {
	if p.Output == nil {
		return
	}
	tags = append(p.Field.DefaultTags, tags...)
	if p.Field.Namespace != "" {
		_, err = fmt.Fprintf(p.Output, "%s %s.%s %s %s\n", metricType, p.Field.Namespace, name, value, strings.Join(tags, ","))
	} else {
		_, err = fmt.Fprintf(p.Output, "%s %s %s %s\n", metricType, name, value, strings.Join(tags, ","))
	}
	return
}

// Flush implemenents stats.Collector.
func (p Printer) Flush() error {
	return nil
}

// Close implemenents stats.Collector.
func (p Printer) Close() error {
	return nil
}
