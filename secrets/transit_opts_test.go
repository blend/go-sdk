package secrets

import (
	"github.com/blend/go-sdk/assert"
	"testing"
)

func TestTransitCreateOptConfig(t *testing.T) {
	a := assert.New(t)
	cfg := CreateTransitKeyConfig{
		Derived: true,
		Convergent: true,
	}

	empty := &CreateTransitKeyConfig{}

	OptCreateTransitConfig(cfg)(empty)

	a.True(empty.Derived)
	a.True(empty.Convergent)
}

func TestTransitUpdateOptConfig(t *testing.T) {
	a := assert.New(t)
	cfg := UpdateTransitKeyConfig{
		Exportable: true,
		AllowPlaintextBackup: true,
	}

	empty := &UpdateTransitKeyConfig{}

	OptUpdateTransitConfig(cfg)(empty)

	a.True(empty.Exportable)
	a.True(empty.AllowPlaintextBackup)
}

func TestTransitCreateOptMisc(t *testing.T) {
	a := assert.New(t)
	empty := &CreateTransitKeyConfig{}

	OptCreateTransitDerived()(empty)
	a.True(empty.Derived)

	empty.Derived = false

	OptCreateTransitConvergent()(empty)
	a.True(empty.Derived)
	a.True(empty.Convergent)

	OptCreateTransitAllowPlaintextBackup()(empty)
	a.True(empty.AllowPlaintextBackup)

	OptCreateTransitExportable()(empty)
	a.True(empty.Exportable)

	err := OptCreateTransitType("not a real type")(empty)
	a.NotNil(err)

	err = OptCreateTransitType(TypeCHACHA20POLY1305)(empty)
	a.Nil(err)
	a.Equal(TypeCHACHA20POLY1305, empty.Type)
}

func TestTransitUpdateOptMisc(t *testing.T) {
	a := assert.New(t)
	empty := &UpdateTransitKeyConfig{}

	OptUpdateTransitDeletionAllowed(true)(empty)
	a.True(*empty.DeletionAllowed)

	OptUpdateTransitDeletionAllowed(false)(empty)
	a.False(*empty.DeletionAllowed)

	OptUpdateTransitAllowPlaintextBackup()(empty)
	a.True(empty.AllowPlaintextBackup)

	OptUpdateTransitExportable()(empty)
	a.True(empty.Exportable)

	OptUpdateTransitMinDecryptionVer(4)(empty)
	a.Equal(4, empty.MinDecryptionVersion)

	OptUpdateTransitMinEncryptionVer(5)(empty)
	a.Equal(5, empty.MinEncryptionVersion)
}
