package vault

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestTransitCreateOptConfig(t *testing.T) {
	a := assert.New(t)
	cfg := CreateTransitKeyConfig{
		Derived:    true,
		Convergent: true,
	}

	empty := &CreateTransitKeyConfig{}

	a.Nil(OptCreateTransitConfig(cfg)(empty))

	a.True(empty.Derived)
	a.True(empty.Convergent)
}

func TestTransitUpdateOptConfig(t *testing.T) {
	a := assert.New(t)
	cfg := UpdateTransitKeyConfig{
		Exportable:           true,
		AllowPlaintextBackup: true,
	}

	empty := &UpdateTransitKeyConfig{}

	a.Nil(OptUpdateTransitConfig(cfg)(empty))

	a.True(empty.Exportable)
	a.True(empty.AllowPlaintextBackup)
}

func TestTransitCreateOptMisc(t *testing.T) {
	a := assert.New(t)
	empty := &CreateTransitKeyConfig{}

	a.Nil(OptCreateTransitDerived()(empty))
	a.True(empty.Derived)

	empty.Derived = false

	a.Nil(OptCreateTransitConvergent()(empty))
	a.True(empty.Derived)
	a.True(empty.Convergent)

	a.Nil(OptCreateTransitAllowPlaintextBackup()(empty))
	a.True(empty.AllowPlaintextBackup)

	a.Nil(OptCreateTransitExportable()(empty))
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

	a.Nil(OptUpdateTransitDeletionAllowed(true)(empty))
	a.True(*empty.DeletionAllowed)

	a.Nil(OptUpdateTransitDeletionAllowed(false)(empty))
	a.False(*empty.DeletionAllowed)

	a.Nil(OptUpdateTransitAllowPlaintextBackup()(empty))
	a.True(empty.AllowPlaintextBackup)

	a.Nil(OptUpdateTransitExportable()(empty))
	a.True(empty.Exportable)

	a.Nil(OptUpdateTransitMinDecryptionVer(4)(empty))
	a.Equal(4, empty.MinDecryptionVersion)

	a.Nil(OptUpdateTransitMinEncryptionVer(5)(empty))
	a.Equal(5, empty.MinEncryptionVersion)
}
