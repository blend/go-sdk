/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package dbutil

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/blend/go-sdk/ex"
)

var (
	// ReservedDatabaseNames are names you cannot use to create a database with.
	ReservedDatabaseNames	= []string{
		"postgres",
		"defaultdb",
		"template0",
		"template1",
	}

	// DatabaseNameMaxLength is the maximum length of a database name.
	DatabaseNameMaxLength	= 63
)

const (
	// ErrDatabaseNameReserved is a validation failure.
	ErrDatabaseNameReserved	ex.Class	= "dbutil; database name is reserved"

	// ErrDatabaseNameEmpty is a validation failure.
	ErrDatabaseNameEmpty	ex.Class	= "dbutil; database name is empty"

	// ErrDatabaseNameInvalidFirstRune is a validation failure.
	ErrDatabaseNameInvalidFirstRune	ex.Class	= "dbutil; database name must start with a letter or underscore"

	// ErrDatabaseNameInvalid is a validation failure.
	ErrDatabaseNameInvalid	ex.Class	= "dbutil; database name must be composed of (in regex form) [a-zA-Z0-9_]"

	// ErrDatabaseNameTooLong is a validation failure.
	ErrDatabaseNameTooLong	ex.Class	= "dbutil; database name must be 63 characters or fewer"
)

// ValidateDatabaseName validates a database name.
func ValidateDatabaseName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ex.New(ErrDatabaseNameEmpty)
	}
	if len(name) > DatabaseNameMaxLength {
		return ex.New(ErrDatabaseNameTooLong)
	}

	firstRune, _ := utf8.DecodeRuneInString(name)
	if !isValidDatabaseNameFirstRune(firstRune) {
		return ex.New(ErrDatabaseNameInvalidFirstRune, ex.OptMessagef("database name: %s", name))
	}

	for _, r := range name {
		if !isValidDatabaseNameRune(r) {
			return ex.New(ErrDatabaseNameInvalid, ex.OptMessagef("database name: %s", name))
		}
	}

	for _, reserved := range ReservedDatabaseNames {
		if strings.EqualFold(reserved, name) {
			return ex.New(ErrDatabaseNameReserved, ex.OptMessagef("database name: %s", name))
		}
	}
	return nil
}

// isValidDatabaseNameFirstRune returns if the rune is valid as a first rune of a database name.
func isValidDatabaseNameFirstRune(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

// isValidDatabaseNameRune is a rune predicate that indicites a rune is a valid database name component.
func isValidDatabaseNameRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}
