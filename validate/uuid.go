/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package validate

import (
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/uuid"
)

// UUID errors
const (
	ErrUUIDRequired		ex.Class	= "uuid should be set"
	ErrUUIDForbidden	ex.Class	= "uuid should not be set"
	ErrUUIDV4		ex.Class	= "uuid should be version 4"
	ErrUUIDVersion		ex.Class	= "uuid should be a given version"
)

// UUID returns a uuid.UUID validator.
func UUID(value *uuid.UUID) UUIDValidators {
	return UUIDValidators{value}
}

// UUIDValidators implements validators for uuid.UUIDs.
type UUIDValidators struct {
	Value *uuid.UUID
}

// Required returns a validator that an uuid is set.
func (u UUIDValidators) Required() Validator {
	return func() error {
		if u.Value == nil {
			return Error(ErrUUIDRequired, nil)
		}
		if len(*u.Value) == 0 {
			return Error(ErrUUIDRequired, nil)
		}
		if u.Value.IsZero() {
			return Error(ErrUUIDRequired, nil)
		}
		return nil
	}
}

// Forbidden returns a validator that an uuid is not set.
func (u UUIDValidators) Forbidden() Validator {
	return func() error {
		if err := u.Required()(); err == nil {
			return Error(ErrUUIDForbidden, nil)
		}
		return nil
	}
}

// IsV4 returns a validator that asserts a uuid.UUID is a V4 uuid.
func (u UUIDValidators) IsV4() Validator {
	return func() error {
		if u.Value == nil {
			return Error(ErrUUIDV4, nil)
		}
		if len(*u.Value) == 0 {
			return Error(ErrUUIDV4, nil)
		}
		if !u.Value.IsV4() {
			return Error(ErrUUIDV4, nil)
		}
		return nil
	}
}

// IsVersion returns a validator that asserts a uuid.UUID is a given version.
func (u UUIDValidators) IsVersion(version byte) Validator {
	return func() error {
		if u.Value == nil {
			return Errorf(ErrUUIDVersion, nil, "version: %x", version)
		}
		if len(*u.Value) == 0 {
			return Errorf(ErrUUIDVersion, nil, "version: %x", version)
		}
		if !u.Value.IsV4() {
			return Errorf(ErrUUIDVersion, nil, "version: %x", version)
		}
		return nil
	}
}
