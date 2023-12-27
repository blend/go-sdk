/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package db

import "strings"

// Dialect is the flavor of sql.
type Dialect string

// Is returns if a dialect equals one of a set of dialects.
func (d Dialect) Is(others ...Dialect) bool {
	for _, other := range others {
		if strings.EqualFold(string(d), string(other)) {
			return true
		}
	}
	return false
}

var (
	// DialectUnknown is an unknown dialect, typically inferred as DialectPostgres.
	DialectUnknown	Dialect	= ""
	// DialectPostgres is the postgres dialect.
	DialectPostgres	Dialect	= "psql"
	// DialectCockroachDB is the crdb dialect.
	DialectCockroachDB	Dialect	= "cockroachdb"
	// DialectRedshift is the redshift dialect.
	DialectRedshift	Dialect	= "redshift"
)
