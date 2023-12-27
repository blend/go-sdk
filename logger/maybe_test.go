/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestMaybeNilLogger(t *testing.T) {
	MaybeInfof(nil, "")
	MaybeInfofContext(context.TODO(), nil, "")
	MaybeInfo(nil, "")
	MaybeInfoContext(context.TODO(), nil, "")
	MaybeDebugf(nil, "")
	MaybeDebugfContext(context.TODO(), nil, "")
	MaybeDebug(nil, "")
	MaybeDebugContext(context.TODO(), nil, "")
	MaybeWarningf(nil, "")
	MaybeWarningfContext(context.TODO(), nil, "")
	MaybeWarning(nil, nil)
	MaybeWarningContext(context.TODO(), nil, nil)
	MaybeErrorf(nil, "")
	MaybeErrorfContext(context.TODO(), nil, "")
	MaybeError(nil, nil)
	MaybeErrorContext(context.TODO(), nil, nil)
	MaybeFatalf(nil, "")
	MaybeFatalfContext(context.TODO(), nil, "")
	MaybeFatal(nil, nil)
	MaybeFatalContext(context.TODO(), nil, nil)
}

func TestMaybeLogger(t *testing.T) {
	assert := assert.New(t)

	log := MustNew(OptAll())
	log.Formatter = NewTextOutputFormatter(
		OptTextNoColor(),
		OptTextHideTimestamp(),
	)

	labelsCtx := func(key, value string) context.Context {
		return WithLabels(context.Background(), Labels{key: value})
	}

	buf := new(bytes.Buffer)
	log.Output = buf
	MaybeInfof(log, "Infof")
	assert.Equal("[info] Infof\n", buf.String())

	buf = new(bytes.Buffer)
	log.Output = buf
	MaybeInfofContext(labelsCtx("a", "b"), log, "Infof")
	assert.Equal("[info] Infof\ta=b\n", buf.String())

	buf = new(bytes.Buffer)
	log.Output = buf
	MaybeInfo(log, "Info")
	assert.Equal("[info] Info\n", buf.String())

	buf = new(bytes.Buffer)
	log.Output = buf
	MaybeInfoContext(labelsCtx("a", "b"), log, "Info")
	assert.Equal("[info] Info\ta=b\n", buf.String())

	buf = new(bytes.Buffer)
	log.Output = buf
	MaybeDebugf(log, "Debugf")
	assert.Equal("[debug] Debugf\n", buf.String())

	buf = new(bytes.Buffer)
	log.Output = buf
	MaybeDebugfContext(labelsCtx("a", "b"), log, "Debugf")
	assert.Equal("[debug] Debugf\ta=b\n", buf.String())

	buf = new(bytes.Buffer)
	log.Output = buf
	MaybeDebug(log, "Debug")
	assert.Equal("[debug] Debug\n", buf.String())

	buf = new(bytes.Buffer)
	log.Output = buf
	MaybeDebugContext(labelsCtx("a", "b"), log, "Debug")
	assert.Equal("[debug] Debug\ta=b\n", buf.String())

	buf = new(bytes.Buffer)
	log.Output = buf
	MaybeWarningf(log, "Warningf")
	assert.Equal("[warning] Warningf\n", buf.String())

	buf = new(bytes.Buffer)
	log.Output = buf
	MaybeWarningfContext(labelsCtx("a", "b"), log, "Warningf")
	assert.Equal("[warning] Warningf\ta=b\n", buf.String())

	buf = new(bytes.Buffer)
	log.Output = buf
	MaybeWarning(log, fmt.Errorf("Warning"))
	assert.Equal("[warning] Warning\n", buf.String())

	buf = new(bytes.Buffer)
	log.Output = buf
	MaybeWarningContext(labelsCtx("a", "b"), log, fmt.Errorf("Warning"))
	assert.Equal("[warning] Warning\ta=b\n", buf.String())

	buf = new(bytes.Buffer)
	log.Output = buf
	MaybeErrorf(log, "Errorf")
	assert.Equal("[error] Errorf\n", buf.String())

	buf = new(bytes.Buffer)
	log.Output = buf
	MaybeErrorfContext(labelsCtx("a", "b"), log, "Errorf")
	assert.Equal("[error] Errorf\ta=b\n", buf.String())

	buf = new(bytes.Buffer)
	log.Output = buf
	MaybeError(log, fmt.Errorf("Error"))
	assert.Equal("[error] Error\n", buf.String())

	buf = new(bytes.Buffer)
	log.Output = buf
	MaybeErrorContext(labelsCtx("a", "b"), log, fmt.Errorf("Error"))
	assert.Equal("[error] Error\ta=b\n", buf.String())

	buf = new(bytes.Buffer)
	log.Output = buf
	MaybeFatalf(log, "Fatalf")
	assert.Equal("[fatal] Fatalf\n", buf.String())

	buf = new(bytes.Buffer)
	log.Output = buf
	MaybeFatalfContext(labelsCtx("a", "b"), log, "Fatalf")
	assert.Equal("[fatal] Fatalf\ta=b\n", buf.String())

	buf = new(bytes.Buffer)
	log.Output = buf
	MaybeFatal(log, fmt.Errorf("Fatal"))
	assert.Equal("[fatal] Fatal\n", buf.String())

	buf = new(bytes.Buffer)
	log.Output = buf
	MaybeFatalContext(labelsCtx("a", "b"), log, fmt.Errorf("Fatal"))
	assert.Equal("[fatal] Fatal\ta=b\n", buf.String())
}
