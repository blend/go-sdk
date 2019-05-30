package logger

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestMaybeNilLogger(t *testing.T) {
	MaybeInfof(nil, "")
	MaybeInfo(nil, "")
	MaybeDebugf(nil, "")
	MaybeDebug(nil, "")
	MaybeWarningf(nil, "")
	MaybeWarning(nil, nil)
	MaybeErrorf(nil, "")
	MaybeError(nil, nil)
	MaybeFatalf(nil, "")
	MaybeFatal(nil, nil)
}

func TestMaybeLogger(t *testing.T) {
	assert := assert.New(t)

	log := MustNew(OptAll())
	log.Formatter = NewTextOutputFormatter(
		OptTextNoColor(),
		OptTextHideTimestamp(),
	)

	buf := new(bytes.Buffer)
	log.Output = buf
	MaybeInfof(log, "Infof")
	assert.Equal("[info] Infof\n", buf.String())

	buf = new(bytes.Buffer)
	log.Output = buf
	MaybeInfo(log, "Info")
	assert.Equal("[info] Info\n", buf.String())

	buf = new(bytes.Buffer)
	log.Output = buf
	MaybeDebugf(log, "Debugf")
	assert.Equal("[debug] Debugf\n", buf.String())

	buf = new(bytes.Buffer)
	log.Output = buf
	MaybeDebug(log, "Debug")
	assert.Equal("[debug] Debug\n", buf.String())

	buf = new(bytes.Buffer)
	log.Output = buf
	MaybeWarningf(log, "Warningf")
	assert.Equal("[warning] Warningf\n", buf.String())

	buf = new(bytes.Buffer)
	log.Output = buf
	MaybeWarning(log, fmt.Errorf("Warning"))
	assert.Equal("[warning] Warning\n", buf.String())

	buf = new(bytes.Buffer)
	log.Output = buf
	MaybeErrorf(log, "Errorf")
	assert.Equal("[error] Errorf\n", buf.String())

	buf = new(bytes.Buffer)
	log.Output = buf
	MaybeError(log, fmt.Errorf("Error"))
	assert.Equal("[error] Error\n", buf.String())

	buf = new(bytes.Buffer)
	log.Output = buf
	MaybeFatalf(log, "Fatalf")
	assert.Equal("[fatal] Fatalf\n", buf.String())

	buf = new(bytes.Buffer)
	log.Output = buf
	MaybeFatal(log, fmt.Errorf("Fatal"))
	assert.Equal("[fatal] Fatal\n", buf.String())
}
