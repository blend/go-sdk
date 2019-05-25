package logger

import (
	"fmt"
	"testing"
)

func TestMaybeNilLogger(t *testing.T) {
	MaybeInfof(nil, "")
	MaybeDebugf(nil, "")
	MaybeWarningf(nil, "")
	MaybeWarning(nil, nil)
	MaybeErrorf(nil, "")
	MaybeError(nil, nil)
	MaybeFatalf(nil, "")
	MaybeFatal(nil, nil)
}

func TestMaybeLogger(t *testing.T) {
	MaybeInfof(MustNew(OptOutput(nil)), "Infof")
	MaybeDebugf(MustNew(OptOutput(nil)), "Debugf")
	MaybeWarningf(MustNew(OptOutput(nil)), "Warningf")
	MaybeWarning(MustNew(OptOutput(nil)), fmt.Errorf("Warning"))
	MaybeErrorf(MustNew(OptOutput(nil)), "Errorf")
	MaybeError(MustNew(OptOutput(nil)), fmt.Errorf("Error"))
	MaybeFatalf(MustNew(OptOutput(nil)), "Fatalf")
	MaybeFatal(MustNew(OptOutput(nil)), fmt.Errorf("Fatal"))
}
