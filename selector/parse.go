/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package selector

// Parse takes a string representing a selector and returns a selector
// object, or an error.
// The input will cause an error if it does not follow this form:
//
//	<selector-syntax>         ::= <requirement> | <requirement> "," <selector-syntax>
//	<requirement>             ::= [!] KEY [ <set-based-restriction> | <exact-match-restriction> ]
//	<set-based-restriction>   ::= "" | <inclusion-exclusion> <value-set>
//	<inclusion-exclusion>     ::= <inclusion> | <exclusion>
//	<exclusion>               ::= "notin"
//	<inclusion>               ::= "in"
//	<value-set>               ::= "(" <values> ")"
//	<values>                  ::= VALUE | VALUE "," <values>
//	<exact-match-restriction> ::= ["="|"=="|"!="] VALUE
//
// KEY is a sequence of one or more characters following: [ DNS_SUBDOMAIN "/" ] DNS_LABEL.
//   - DNS_SUBDOMAIN is a sequence of one or more characters ([a-z0-9-.]), and must start and end with an alphanumeric ([a-z0-9]).
//     Symbol characters ('.', '-') cannot repeat. Max length is 253 characters.
//   - DNS_LABEL is a sequence of one or more characters ([A-Za-z0-9_-.]). Max length is 63 characters.
//
// VALUE is a sequence of zero or more characters ([A-Za-z0-9_-.]). Values must start and end with an alphanumeric ([a-z0-9A-Z]). Max length is 63 characters.
// Delimiter is white space: (' ', '\t')
//
// Example of valid syntax:
//
//	"x in (foo,,baz),y,z notin ()"
//
// Note:
//
//	(1) Inclusion - " in " - denotes that the KEY exists and is equal to any of the
//	    VALUEs in its requirement
//	(2) Exclusion - " notin " - denotes that the KEY is not equal to any
//	    of the VALUEs in its requirement or does not exist
//	(3) The empty string is a valid VALUE
//	(4) A requirement with just a KEY - as in "y" above - denotes that
//	    the KEY exists and can be any VALUE.
//	(5) A requirement with just !KEY requires that the KEY not exist.
func Parse(query string, opts ...Option) (Selector, error) {
	p := &Parser{s: query}
	for _, opt := range opts {
		opt(p)
	}
	return p.Parse()
}

// MustParse parses the selector but will panic if there is an issue.
func MustParse(query string, opts ...Option) Selector {
	sel, err := Parse(query, opts...)
	if err != nil {
		panic(err)
	}
	return sel
}
