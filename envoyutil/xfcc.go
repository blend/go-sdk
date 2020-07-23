package envoyutil

import (
	"fmt"
	"strings"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/stringutil"
)

// XFCC represents a proxy header containing certificate information for the client
// that is sending the request to the proxy.
// See https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_conn_man/headers#x-forwarded-client-cert
type XFCC []XFCCElement

// XFCCElement is an element in an XFCC header (see `XFCC`).
//
// NOTE: This is an intentionally limited coverage of the fields available in the XFCC
//       header.
type XFCCElement struct {
	By  string
	URI string
}

// maybeQuoted quotes a string value that may need to be quoted to be part of an
// XFCC header. It will use `%q` formatting to quote the value if it contains any
// of `,` (comma), `;` (semi-colon), `=` (equals) or `"` (double quote).
func maybeQuoted(value string) string {
	if strings.ContainsAny(value, `,;="`) {
		return fmt.Sprintf("%q", value)
	}
	return value
}

// String converts the parsed XFCC element **back** to a string. This is intended
// for debugging purposes and is not particularly
func (xe XFCCElement) String() string {
	parts := []string{}
	if xe.By != "" {
		parts = append(parts, fmt.Sprintf("By=%s", maybeQuoted(xe.By)))
	}
	if xe.URI != "" {
		parts = append(parts, fmt.Sprintf("URI=%s", maybeQuoted(xe.URI)))
	}
	return strings.Join(parts, ";")
}

const (
	// HeaderXFCC is the header key for forwarded client cert
	HeaderXFCC = "x-forwarded-client-cert"
)

const (
	// ErrXFCCParsing is the class of error returned when parsing XFCC fails
	ErrXFCCParsing = ex.Class("Error Parsing X-Forwarded-Client-Cert")
)

type parseXFCCState int

const (
	parseXFCCKey parseXFCCState = iota
	parseXFCCValue
)

// ParseXFCC parses the XFCC header
func ParseXFCC(header string) (XFCC, error) {
	xfcc := XFCC{}
	elements := stringutil.SplitCSV(header)
	for _, element := range elements {
		ele, err := ParseXFCCElement(element)
		if err != nil {
			return XFCC{}, err
		}
		xfcc = append(xfcc, ele)
	}
	return xfcc, nil
}

// ParseXFCCElement parses an element out of the given string. An error is returned if the parser
// encounters a key not in the valid list or the string is malformed
func ParseXFCCElement(element string) (XFCCElement, error) {
	state := parseXFCCKey
	ele := XFCCElement{}
	key := ""
	valueStart := -1
	valueEnd := -1
	for i, char := range element {
		switch state {
		case parseXFCCKey:
			if char == '=' {
				state = parseXFCCValue
			} else {
				key += string(char)
			}
		case parseXFCCValue:
			if char == ';' {
				if len(key) == 0 || valueStart == -1 {
					return XFCCElement{}, ex.New(ErrXFCCParsing).WithMessage("Key or Value missing")
				}
				err := fillXFCCKeyValue(key, element, valueStart, valueEnd, &ele)
				if err != nil {
					return XFCCElement{}, err
				}

				key = ""
				valueStart = -1
				valueEnd = -1
				state = parseXFCCKey
			} else {
				if valueStart == -1 {
					valueStart = i
				}
				valueEnd = i
			}
		}
	}

	if len(key) > 0 && valueStart != -1 {
		return ele, fillXFCCKeyValue(key, element, valueStart, valueEnd, &ele)
	}

	if len(key) > 0 || valueStart != -1 {
		return XFCCElement{}, ex.New(ErrXFCCParsing).WithMessage("Key or value found but not both")
	}

	return ele, nil
}

func fillXFCCKeyValue(key, element string, valueStart, valueEnd int, ele *XFCCElement) (err error) {
	key = strings.ToLower(key)
	switch key {
	case "cert", "chain", "dns", "hash", "subject":
		return nil
	case "by":
		if ele.By != "" {
			return ex.New(ErrXFCCParsing).WithMessagef("Key already encountered %q", key)
		}
		// NOTE: This can panic at runtime if `valueStart` and / or `valueEnd` are malformed.
		//       The assumption here is that the "valid range" invariant for these inputs is maintained
		//       elsewhere, in `ParseXFCCElement()`.
		ele.By = element[valueStart : valueEnd+1]
	case "uri":
		if ele.URI != "" {
			return ex.New(ErrXFCCParsing).WithMessagef("Key already encountered %q", key)
		}
		// NOTE: This can panic at runtime if `valueStart` and / or `valueEnd` are malformed.
		//       The assumption here is that the "valid range" invariant for these inputs is maintained
		//       elsewhere, in `ParseXFCCElement()`.
		ele.URI = element[valueStart : valueEnd+1]
	default:
		return ex.New(ErrXFCCParsing).WithMessagef("Unknown key %q", key)
	}
	return nil
}
