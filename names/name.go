/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package names

// Name is a structured/parsed name.
type Name struct {
	Salutation	string
	FirstName	string
	MiddleName	string
	LastName	string
	Suffix		string
}

// String returns the string representation of a name.
func (n Name) String() string {
	fullName := ""

	if n.Salutation != "" {
		fullName += n.Salutation
	}

	if n.FirstName != "" {
		if fullName != "" {
			fullName += " "
		}
		fullName += n.FirstName
	}

	if n.MiddleName != "" {
		if fullName != "" {
			fullName += " "
		}
		fullName += n.MiddleName
	}

	if n.LastName != "" {
		if fullName != "" {
			fullName += " "
		}
		fullName += n.LastName
	}
	if n.Suffix != "" {
		if fullName != "" {
			fullName += " "
		}
		fullName += n.Suffix
	}

	return fullName
}
