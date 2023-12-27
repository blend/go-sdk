/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

import (
	"context"
	"encoding/json"
	"io"

	"github.com/blend/go-sdk/bufferutil"
)

var (
	_ WriteFormatter = (*JSONOutputFormatter)(nil)
)

// NewJSONOutputFormatter returns a new json event formatter.
func NewJSONOutputFormatter(options ...JSONOutputFormatterOption) *JSONOutputFormatter {
	jf := &JSONOutputFormatter{
		BufferPool: bufferutil.NewPool(DefaultBufferPoolSize),
	}

	for _, option := range options {
		option(jf)
	}
	return jf
}

// JSONOutputFormatterOption is an option for json formatters.
type JSONOutputFormatterOption func(*JSONOutputFormatter)

// OptJSONConfig sets a json formatter from a config.
func OptJSONConfig(cfg JSONConfig) JSONOutputFormatterOption {
	return func(jf *JSONOutputFormatter) {
		jf.Pretty = cfg.Pretty
		jf.PrettyIndent = cfg.PrettyIndentOrDefault()
		jf.PrettyPrefix = cfg.PrettyPrefixOrDefault()
	}
}

// OptJSONPretty sets the json output formatter to indent output.
func OptJSONPretty() JSONOutputFormatterOption {
	return func(jso *JSONOutputFormatter) { jso.Pretty = true }
}

// OptJSONPrettyPrefix sets the json output formatter to indent output.
func OptJSONPrettyPrefix(prettyPrefix string) JSONOutputFormatterOption {
	return func(jso *JSONOutputFormatter) { jso.PrettyPrefix = prettyPrefix }
}

// OptJSONPrettyIndent sets the json output formatter to indent output.
func OptJSONPrettyIndent(prettyIndent string) JSONOutputFormatterOption {
	return func(jso *JSONOutputFormatter) { jso.PrettyIndent = prettyIndent }
}

// JSONOutputFormatter is a json output formatter.
type JSONOutputFormatter struct {
	BufferPool   *bufferutil.Pool
	Pretty       bool
	PrettyPrefix string
	PrettyIndent string
}

// PrettyPrefixOrDefault returns the pretty prefix or a default.
func (jw JSONOutputFormatter) PrettyPrefixOrDefault() string {
	if jw.PrettyPrefix != "" {
		return jw.PrettyPrefix
	}
	return ""
}

// PrettyIndentOrDefault returns the pretty indent or a default.
func (jw JSONOutputFormatter) PrettyIndentOrDefault() string {
	if jw.PrettyIndent != "" {
		return jw.PrettyIndent
	}
	return "\t"
}

// WriteFormat writes the event to the given output.
func (jw JSONOutputFormatter) WriteFormat(ctx context.Context, output io.Writer, e Event) error {
	buffer := jw.BufferPool.Get()
	defer jw.BufferPool.Put(buffer)

	encoder := json.NewEncoder(buffer)
	if jw.Pretty {
		encoder.SetIndent(jw.PrettyPrefixOrDefault(), jw.PrettyIndentOrDefault())
	}
	if decomposer, ok := e.(JSONWritable); ok {
		fields := jw.CombineFields(jw.GetScopeFields(ctx, e), decomposer.Decompose())
		if err := encoder.Encode(fields); err != nil {
			return err
		}
	} else if err := encoder.Encode(e); err != nil {
		return err
	}
	_, err := io.Copy(output, buffer)
	return err
}

// CombineFields combines a variadic set of fields.
func (jw JSONOutputFormatter) CombineFields(fields ...map[string]interface{}) map[string]interface{} {
	if len(fields) == 0 {
		return nil
	}
	output := make(map[string]interface{})
	for _, set := range fields {
		if set == nil {
			continue
		}
		for key, value := range set {
			output[key] = value
		}
	}
	return output
}

// GetScopeFields gets scope fields from a context.
func (jw JSONOutputFormatter) GetScopeFields(ctx context.Context, e Event) map[string]interface{} {
	output := map[string]interface{}{
		FieldFlag:       e.GetFlag(),
		FieldTimestamp:  GetEventTimestamp(ctx, e),
		FieldRestricted: GetRestricted(ctx),
	}
	if path := GetPath(ctx); len(path) > 0 {
		output[FieldScopePath] = path
	}
	if labels := GetLabels(ctx); len(labels) > 0 {
		output[FieldLabels] = labels
	}
	if annotations := GetAnnotations(ctx); len(annotations) > 0 {
		output[FieldAnnotations] = annotations
	}
	return output
}
