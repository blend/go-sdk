package names

// Name is a structured/parsed name.
type Name struct {
	Salutation string
	FirstName  string
	MiddleName string
	LastName   string
	Suffix     string
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
