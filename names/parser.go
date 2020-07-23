package names

import (
	"strings"
	"unicode"
)

var validSuffixes = []string{
	"I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX", "X", "XI", "XII", "XIII", "XIV", "XV", "XVI", "XVII", "XVIII", "XIX", "XX",
	"Senior", "Junior", "Jr", "Sr",
	"PhD", "APR", "RPh", "PE", "MD", "MA", "DMD", "CME",
}

var compoundLastNames = []string{
	"vere", "von", "van", "de", "del", "della", "di", "da", "pietro",
	"vanden", "du", "st.", "st", "la", "lo", "ter", "bin", "ibn",
}

// Parse parses a string into a name.
func Parse(input string) (name Name) {
	fullName := strings.TrimSpace(input)

	rawNameParts := strings.Split(fullName, " ")

	nameParts := []string{}

	lastName := ""
	firstName := ""
	initials := ""
	for _, part := range rawNameParts {
		if !strings.Contains(part, "(") {
			nameParts = append(nameParts, part)
		}
	}

	numWords := len(nameParts)
	salutation := processSalutation(nameParts[0])
	suffix := processSuffix(nameParts[len(nameParts)-1])

	start := 0
	if salutation != "" {
		start = 1
	}

	end := numWords
	if suffix != "" {
		end = numWords - 1
	}

	i := 0
	for i = start; i < (end - 1); i++ {
		word := nameParts[i]
		if isCompoundLastName(word) && i != start {
			break
		}
		if isMiddleName(word) {
			if i == start {
				if isMiddleName(nameParts[i+1]) {
					firstName = firstName + " " + strings.ToUpper(word)
				} else {
					initials = initials + " " + strings.ToUpper(word)
				}
			} else {
				initials = initials + " " + strings.ToUpper(word)
			}
		} else {
			firstName = firstName + " " + fixCase(word)
		}
	}

	if (end - start) > 1 {
		for j := i; j < end; j++ {
			lastName = lastName + " " + fixCase(nameParts[j])
		}
	} else if i < len(nameParts) {
		firstName = fixCase(nameParts[i])
	}

	name.Salutation = salutation
	name.FirstName = strings.TrimSpace(firstName)
	name.MiddleName = strings.TrimSpace(initials)
	name.LastName = strings.TrimSpace(lastName)
	name.Suffix = suffix

	return name
}

func processSalutation(input string) string {
	word := cleanString(input)

	switch word {
	case "mr", "master", "mister":
		return "Mr."
	case "mrs", "misses":
		return "Mrs."
	case "ms", "miss":
		return "Ms."
	case "dr":
		return "Dr."
	case "rev":
		return "Rev."
	case "fr":
		return "Fr."
	}

	return ""
}

func processSuffix(input string) string {
	word := cleanString(input)
	return getByLower(validSuffixes, word)
}

func isCompoundLastName(input string) bool {
	word := cleanString(input)
	exists := containsLower(compoundLastNames, word)
	return exists
}

func isMiddleName(input string) bool {
	word := cleanString(input)
	return len(word) == 1
}

func uppercaseFirstAll(input string, seperator string) string {
	words := []string{}
	parts := strings.Split(input, seperator)
	for _, thisWord := range parts {
		toAppend := ""
		switch {
		case isCompoundLastName(strings.ToLower(thisWord)):
			// preserve first letter case, but to lower the rest for compound last names
			if unicode.IsUpper([]rune(thisWord)[0]) {
				toAppend = strings.Title(strings.ToLower(thisWord))
			} else {
				toAppend = strings.ToLower(thisWord)
			}
		case isCamelCase(thisWord):
			// Preserve case for Camel-cased strings
			toAppend = thisWord
		default:
			// For everything else, force to title case
			toAppend = upperCaseFirst(strings.ToLower(thisWord))
		}
		words = append(words, toAppend)
	}
	return strings.Join(words, seperator)
}

func upperCaseFirst(input string) string {
	return strings.Title(strings.ToLower(input))
}

func fixCase(input string) string {
	word := uppercaseFirstAll(input, "-")
	word = uppercaseFirstAll(word, ".")
	return word
}

func cleanString(input string) string {
	return strings.ToLower(strings.Replace(input, ".", "", -1))
}

// isCamelCase returns if a string is CamelCased.
// CamelCased in this sense is if a string has both inner-upper and lower characters.
func isCamelCase(input string) bool {
	hasLowers := false
	hasInnerUppers := false

	for i, c := range input {
		if i != 0 && unicode.IsUpper(c) {
			hasInnerUppers = true
		}
		if unicode.IsLower(c) {
			hasLowers = true
		}
	}

	return hasLowers && hasInnerUppers
}

// containsLower returns true if the `elem` is in the StringArray, false otherwise.
func containsLower(values []string, elem string) bool {
	for _, arrayElem := range values {
		if strings.ToLower(arrayElem) == elem {
			return true
		}
	}
	return false
}

// getByLower returns an element from the array that matches the input.
func getByLower(values []string, elem string) string {
	for _, arrayElem := range values {
		if strings.ToLower(arrayElem) == elem {
			return arrayElem
		}
	}
	return ""
}
