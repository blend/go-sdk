/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package db

import (
	"net"

	"github.com/blend/go-sdk/ex"
)

const (
	// ErrDestinationNotStruct is an exception class.
	ErrDestinationNotStruct ex.Class = "db: destination object is not a struct"
	// ErrConfigUnset is an exception class.
	ErrConfigUnset ex.Class = "db: config is unset"
	// ErrUnsafeSSLMode is an error indicating unsafe ssl mode in production.
	ErrUnsafeSSLMode ex.Class = "db: unsafe ssl mode in prodlike environment"
	// ErrUsernameUnset is an error indicating there is no username set in a prodlike environment.
	ErrUsernameUnset ex.Class = "db: username is unset in prodlike environment"
	// ErrPasswordUnset is an error indicating there is no password set in a prodlike environment.
	ErrPasswordUnset ex.Class = "db: password is unset in prodlike environment"
	// ErrDurationConversion is the error returned when a duration cannot be
	// converted to multiple of some base (e.g. milliseconds or seconds)
	// without round off.
	ErrDurationConversion ex.Class = "db: cannot convert duration"
	// ErrConnectionAlreadyOpen is an error indicating the db connection was already opened.
	ErrConnectionAlreadyOpen ex.Class = "db: the connection is already opened"
	// ErrConnectionClosed is an error indicating the db connection hasn't been opened.
	ErrConnectionClosed ex.Class = "db: the connection is closed, or is being used before opened"
	// ErrPlanCacheUnset is an error indicating the statement cache is unset.
	ErrPlanCacheUnset ex.Class = "db: the plan cache is unset"
	// ErrPlanCacheKeyUnset is an error indicating the plan cache key is unset.
	ErrPlanCacheKeyUnset ex.Class = "db: the plan cache key is unset"
	// ErrCollectionNotSlice is an error returned by OutMany if the destination is not a slice.
	ErrCollectionNotSlice ex.Class = "db: outmany destination collection is not a slice"
	// ErrInvalidIDs is an error returned by Get if the ids aren't provided.
	ErrInvalidIDs ex.Class = "db: invalid `ids` parameter"
	// ErrNoPrimaryKey is an error returned by a number of operations that depend on a primary key.
	ErrNoPrimaryKey ex.Class = "db: no primary key on object"
	// ErrRowsNotColumnsProvider is returned by `PopulateByName` if you do not pass in `sql.Rows` as the scanner.
	ErrRowsNotColumnsProvider ex.Class = "db: rows is not a columns provider"
	// ErrTooManyRows is returned by Out if there is more than one row returned by the query
	ErrTooManyRows ex.Class = "db: too many rows returned to map to single object"

	// ErrNetwork is a grouped error for network issues.
	ErrNetwork ex.Class = "db: network error"
)

// IsConfigUnset returns if the error is an `ErrConfigUnset`.
func IsConfigUnset(err error) bool {
	return ex.Is(err, ErrConfigUnset)
}

// IsUnsafeSSLMode returns if an error is an `ErrUnsafeSSLMode`.
func IsUnsafeSSLMode(err error) bool {
	return ex.Is(err, ErrUnsafeSSLMode)
}

// IsUsernameUnset returns if an error is an `ErrUsernameUnset`.
func IsUsernameUnset(err error) bool {
	return ex.Is(err, ErrUsernameUnset)
}

// IsPasswordUnset returns if an error is an `ErrPasswordUnset`.
func IsPasswordUnset(err error) bool {
	return ex.Is(err, ErrPasswordUnset)
}

// IsDurationConversion returns if an error is an `ErrDurationConversion`.
func IsDurationConversion(err error) bool {
	return ex.Is(err, ErrDurationConversion)
}

// IsConnectionClosed returns if the error is an `ErrConnectionClosed`.
func IsConnectionClosed(err error) bool {
	return ex.Is(err, ErrConnectionClosed)
}

// IsPlanCacheUnset returns if the error is an `ErrConnectionClosed`.
func IsPlanCacheUnset(err error) bool {
	return ex.Is(err, ErrPlanCacheUnset)
}

// IsPlanCacheKeyUnset returns if the error is an `ErrPlanCacheKeyUnset`.
func IsPlanCacheKeyUnset(err error) bool {
	return ex.Is(err, ErrPlanCacheKeyUnset)
}

// Error returns a new exception by parsing (potentially)
// a driver error into relevant pieces.
func Error(err error, options ...ex.Option) error {
	if err == nil {
		return nil
	}
	if _, ok := err.(*net.OpError); ok {
		return ex.New(ErrNetwork, ex.OptInner(err))
	}
	return ex.New(err, options...)
}
