package envoyutil

import (
	"encoding/hex"
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
type XFCCElement struct {
	// By contains Subject Alternative Name (URI type) of the current proxy's
	// certificate.	It is present here as a `string` and can be parsed to a
	// `url.URL` if desired.
	By string
	// Hash contains the SHA 256 digest of the current client certificate; this
	// is a string of 64 hexadecimal characters. This can be converted to the raw
	// bytes underlying the hex string via `xe.DecodeHash()`.
	Hash string
	// Cert contains the entire client certificate in URL encoded PEM format.
	// It is present here as a `string` and can be parsed to an `x509.Certificate`
	// if desired.
	Cert string
	// Chain contains entire client certificate chain (including the leaf certificate)
	// in URL encoded PEM format. It is present here as a `string` and can be parsed
	// to a `[]x509.Certificate` if desired.
	Chain string
	// Subject contains the `Subject` field of the current client certificate.
	Subject string
	// URI contains the URI SAN of the current client certificate (assumes only
	// one URI SAN).
	URI string
	// DNS contains the DNS SANs of the current client certificate. A client
	// certificate may contain multiple DNS SANs, each will be a separate
	// key-value pair in the XFCC element.
	DNS []string
}

// DecodeHash decodes the `Hash` element from a hex string to raw bytes.
func (xe XFCCElement) DecodeHash() ([]byte, error) {
	return hex.DecodeString(xe.Hash)
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
	value := []rune{}
	for _, char := range element {
		switch state {
		case parseXFCCKey:
			if char == '=' {
				state = parseXFCCValue
			} else {
				key += string(char)
			}
		case parseXFCCValue:
			if char == ';' {
				if len(key) == 0 || len(value) == 0 {
					return XFCCElement{}, ex.New(ErrXFCCParsing).WithMessage("Key or Value missing")
				}
				err := fillXFCCKeyValue(key, element, value, &ele)
				if err != nil {
					return XFCCElement{}, err
				}

				key = ""
				value = []rune{}
				state = parseXFCCKey
			} else {
				value = append(value, char)
			}
		}
	}

	if len(key) > 0 && len(value) > 0 {
		return ele, fillXFCCKeyValue(key, element, value, &ele)
	}

	if len(key) > 0 || len(value) > 0 {
		return XFCCElement{}, ex.New(ErrXFCCParsing).WithMessage("Key or value found but not both")
	}

	return ele, nil
}

func fillXFCCKeyValue(key, element string, value []rune, ele *XFCCElement) (err error) {
	key = strings.ToLower(key)
	switch key {
	case "by":
		if ele.By != "" {
			return ex.New(ErrXFCCParsing).WithMessagef("Key already encountered %q", key)
		}
		ele.By = string(value)
	case "hash":
		if len(ele.Hash) > 0 {
			return ex.New(ErrXFCCParsing).WithMessagef("Key already encountered %q", key)
		}
		ele.Hash = string(value)
	case "cert":
		if len(ele.Cert) > 0 {
			return ex.New(ErrXFCCParsing).WithMessagef("Key already encountered %q", key)
		}
		ele.Cert = string(value)
	case "chain":
		return nil
	case "subject":
		return nil
	case "uri":
		if ele.URI != "" {
			return ex.New(ErrXFCCParsing).WithMessagef("Key already encountered %q", key)
		}
		ele.URI = string(value)
	case "dns":
		return nil
	default:
		return ex.New(ErrXFCCParsing).WithMessagef("Unknown key %q", key)
	}
	return nil
}
