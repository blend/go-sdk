package envoyutil

import (
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"

	"github.com/blend/go-sdk/certutil"
	"github.com/blend/go-sdk/ex"
)

// XFCC represents a proxy header containing certificate information for the client
// that is sending the request to the proxy.
// See https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_conn_man/headers#x-forwarded-client-cert
type XFCC []XFCCElement

// XFCCElement is an element in an XFCC header (see `XFCC`).
type XFCCElement struct {
	// By contains Subject Alternative Name (URI type) of the current proxy's
	// certificate.	This can be decoded as a `*url.URL` via `xe.DecodeBy()`.
	By string
	// Hash contains the SHA 256 digest of the current client certificate; this
	// is a string of 64 hexadecimal characters. This can be converted to the raw
	// bytes underlying the hex string via `xe.DecodeHash()`.
	Hash string
	// Cert contains the entire client certificate in URL encoded PEM format.
	// This can be decoded as a `*x509.Certificate` via `xe.DecodeCert()`.
	Cert string
	// Chain contains entire client certificate chain (including the leaf certificate)
	// in URL encoded PEM format. This can be decoded as a `[]*x509.Certificate` via
	// `xe.DecodeChain()`.
	Chain string
	// Subject contains the `Subject` field of the current client certificate.
	Subject string
	// URI contains the URI SAN of the current client certificate (assumes only
	// one URI SAN). This can be decoded as a `*url.URL` via `xe.DecodeURI()`.
	URI string
	// DNS contains the DNS SANs of the current client certificate. A client
	// certificate may contain multiple DNS SANs, each will be a separate
	// key-value pair in the XFCC element.
	DNS []string
}

// DecodeBy decodes the `By` element from a URI string to a `*url.URL`.
func (xe XFCCElement) DecodeBy() (*url.URL, error) {
	u, err := url.Parse(xe.By)
	if err != nil {
		return nil, ex.New(err)
	}

	return u, nil
}

// DecodeHash decodes the `Hash` element from a hex string to raw bytes.
func (xe XFCCElement) DecodeHash() ([]byte, error) {
	bs, err := hex.DecodeString(xe.Hash)
	if err != nil {
		return nil, ex.New(err)
	}

	return bs, nil
}

// DecodeCert decodes the `Cert` element from a URL encoded PEM to a
// single `x509.Certificate`.
func (xe XFCCElement) DecodeCert() (*x509.Certificate, error) {
	if xe.Cert == "" {
		return nil, nil
	}

	value, err := url.QueryUnescape(xe.Cert)
	if err != nil {
		return nil, ex.New(err)
	}

	parsed, err := certutil.ParseCertPEM([]byte(value))
	if err != nil {
		return nil, ex.New(err)
	}

	if len(parsed) != 1 {
		err = ex.New(
			ErrXFCCParsing,
			ex.OptMessagef("Incorrect number of certificates; expected 1 got %d", len(parsed)),
		)
		return nil, err
	}

	return parsed[0], nil
}

// DecodeChain decodes the `Chain` element from a URL encoded PEM to a
// `[]x509.Certificate`.
func (xe XFCCElement) DecodeChain() ([]*x509.Certificate, error) {
	if xe.Chain == "" {
		return nil, nil
	}

	value, err := url.QueryUnescape(xe.Chain)
	if err != nil {
		return nil, ex.New(err)
	}

	parsed, err := certutil.ParseCertPEM([]byte(value))
	if err != nil {
		return nil, ex.New(err)
	}

	return parsed, nil

}

// DecodeURI decodes the `URI` element from a URI string to a `*url.URL`.
func (xe XFCCElement) DecodeURI() (*url.URL, error) {
	u, err := url.Parse(xe.URI)
	if err != nil {
		return nil, ex.New(err)
	}

	return u, nil
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
	if xe.Hash != "" {
		parts = append(parts, fmt.Sprintf("Hash=%s", maybeQuoted(xe.Hash)))
	}
	if xe.Cert != "" {
		parts = append(parts, fmt.Sprintf("Cert=%s", maybeQuoted(xe.Cert)))
	}
	if xe.Chain != "" {
		parts = append(parts, fmt.Sprintf("Chain=%s", maybeQuoted(xe.Chain)))
	}
	if xe.Subject != "" {
		parts = append(parts, fmt.Sprintf("Subject=%q", xe.Subject))
	}
	if xe.URI != "" {
		parts = append(parts, fmt.Sprintf("URI=%s", maybeQuoted(xe.URI)))
	}
	for _, dnsSAN := range xe.DNS {
		parts = append(parts, fmt.Sprintf("DNS=%s", maybeQuoted(dnsSAN)))
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

	// initialValueCapacity is the capacity used for a value in a key-value
	// pair from an XFCC header.
	initialValueCapacity = 8
)

type parseXFCCState int

const (
	parseXFCCKey parseXFCCState = iota
	parseXFCCValueStart
	parseXFCCValue
	parseXFCCValueQuoted
)

// ParseXFCC parses the XFCC header
func ParseXFCC(header string) (XFCC, error) {
	xfcc := XFCC{}
	elements := strings.Split(header, ",")
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
	asRunes := []rune(element)
	value := make([]rune, 0, initialValueCapacity)
	i := 0
	for i < len(asRunes) {
		char := asRunes[i]
		switch state {
		case parseXFCCKey:
			if char == '=' {
				state = parseXFCCValueStart
			} else {
				key += string(char)
			}
		case parseXFCCValueStart:
			if char == '"' {
				state = parseXFCCValueQuoted
			} else {
				value = append(value, char)
				state = parseXFCCValue
			}
		case parseXFCCValue:
			if char == ';' {
				if len(key) == 0 || len(value) == 0 {
					return XFCCElement{}, ex.New(ErrXFCCParsing).WithMessage("Key or Value missing")
				}
				err := fillXFCCKeyValue(key, value, &ele)
				if err != nil {
					return XFCCElement{}, err
				}

				key = ""
				value = make([]rune, 0, initialValueCapacity)
				state = parseXFCCKey
			} else {
				value = append(value, char)
			}
		case parseXFCCValueQuoted:
			if char == '\\' {
				nextIndex := i + 1
				if nextIndex < len(asRunes) && asRunes[nextIndex] == '"' {
					// Consume two characters at once here (since we have an
					// escaped quote).
					value = append(value, '"')
					i = nextIndex
				} else {
					value = append(value, char)
				}
			} else if char == '"' {
				// Since the **escaped quote** case `\"` has already been
				// covered, this case should only occur in the closing quote
				// case.
				nextIndex := i + 1
				if nextIndex < len(asRunes) {
					if asRunes[nextIndex] != ';' {
						return ele, ex.New(ErrXFCCParsing).WithMessage("Closing quote not followed by `;`.")
					}

					// Consume two characters at once here (since we have an
					// closing quote).
					i = nextIndex

					if len(key) == 0 {
						// Quoted values, e.g. `""`, are allowed to be empty.
						return XFCCElement{}, ex.New(ErrXFCCParsing).WithMessage("Key missing")
					}
					err := fillXFCCKeyValue(key, value, &ele)
					if err != nil {
						return XFCCElement{}, err
					}

					key = ""
					value = []rune{}
					state = parseXFCCKey
				} else {
					// NOTE: If `nextIndex >= len(asRunes)` then we are at the end,
					//       which is a no-op here.
					state = parseXFCCKey
				}
			} else {
				value = append(value, char)
			}
		}

		// Increment `i` for the next iteration. (Note that branches of the `switch`
		// statement may have already incremented `i` as well.)
		i++
	}

	if len(key) > 0 && len(value) > 0 {
		return ele, fillXFCCKeyValue(key, value, &ele)
	}

	if len(key) > 0 || len(value) > 0 {
		return XFCCElement{}, ex.New(ErrXFCCParsing).WithMessage("Key or value found but not both")
	}

	return ele, nil
}

func fillXFCCKeyValue(key string, value []rune, ele *XFCCElement) (err error) {
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
		if len(ele.Chain) > 0 {
			return ex.New(ErrXFCCParsing).WithMessagef("Key already encountered %q", key)
		}
		ele.Chain = string(value)
	case "subject":
		if len(ele.Subject) > 0 {
			return ex.New(ErrXFCCParsing).WithMessagef("Key already encountered %q", key)
		}
		ele.Subject = string(value)
	case "uri":
		if ele.URI != "" {
			return ex.New(ErrXFCCParsing).WithMessagef("Key already encountered %q", key)
		}
		ele.URI = string(value)
	case "dns":
		ele.DNS = append(ele.DNS, string(value))
	default:
		return ex.New(ErrXFCCParsing).WithMessagef("Unknown key %q", key)
	}
	return nil
}
