/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package webutil

import (
	"strings"
)

// IsValidMethod returns if a http method is valid.
/*
	Method         = "OPTIONS"                ; Section 9.2
					| "GET"                    ; Section 9.3
					| "HEAD"                   ; Section 9.4
					| "POST"                   ; Section 9.5
					| "PUT"                    ; Section 9.6
					| "DELETE"                 ; Section 9.7
					| "TRACE"                  ; Section 9.8
					| "CONNECT"                ; Section 9.9
					| extension-method
	extension-method = token
		token          = 1*<any CHAR except CTLs or separators>
*/
func IsValidMethod(method string) bool {
	return len(method) > 0 && strings.IndexFunc(method, isNotMethodToken) == -1
}

func isNotMethodToken(r rune) bool {
	return !IsHTTPTokenRune(r)
}

// IsHTTPTokenRune returns if a rune is in the http valid rune table.
func IsHTTPTokenRune(r rune) bool {
	i := int(r)
	return i < len(isTokenTable) && isTokenTable[i]
}

var isTokenTable = [127]bool{
	'!':	true,
	'#':	true,
	'$':	true,
	'%':	true,
	'&':	true,
	'\'':	true,
	'*':	true,
	'+':	true,
	'-':	true,
	'.':	true,
	'0':	true,
	'1':	true,
	'2':	true,
	'3':	true,
	'4':	true,
	'5':	true,
	'6':	true,
	'7':	true,
	'8':	true,
	'9':	true,
	'A':	true,
	'B':	true,
	'C':	true,
	'D':	true,
	'E':	true,
	'F':	true,
	'G':	true,
	'H':	true,
	'I':	true,
	'J':	true,
	'K':	true,
	'L':	true,
	'M':	true,
	'N':	true,
	'O':	true,
	'P':	true,
	'Q':	true,
	'R':	true,
	'S':	true,
	'T':	true,
	'U':	true,
	'W':	true,
	'V':	true,
	'X':	true,
	'Y':	true,
	'Z':	true,
	'^':	true,
	'_':	true,
	'`':	true,
	'a':	true,
	'b':	true,
	'c':	true,
	'd':	true,
	'e':	true,
	'f':	true,
	'g':	true,
	'h':	true,
	'i':	true,
	'j':	true,
	'k':	true,
	'l':	true,
	'm':	true,
	'n':	true,
	'o':	true,
	'p':	true,
	'q':	true,
	'r':	true,
	's':	true,
	't':	true,
	'u':	true,
	'v':	true,
	'w':	true,
	'x':	true,
	'y':	true,
	'z':	true,
	'|':	true,
	'~':	true,
}
