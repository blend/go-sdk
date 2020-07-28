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

	// initialValueCapacity is the capacity used for a key in a key-value
	// pair from an XFCC header.
	initialKeyCapacity = 4
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

// xfccParser holds state while an XFCC header is being parsed.
type xfccParser struct {
	Header  []rune
	Index   int
	State   parseXFCCState
	Key     []rune
	Value   []rune
	Element XFCCElement
	Parsed  XFCC
}

// ParseXFCC parses the XFCC header.
func ParseXFCC(header string) (XFCC, error) {
	if header == "" {
		return XFCC{}, nil
	}

	xp := &xfccParser{
		Header: []rune(header),
		Index:  0,
		State:  parseXFCCKey,
		Key:    make([]rune, 0, initialKeyCapacity),
		Value:  make([]rune, 0, initialValueCapacity),
	}
	lastCharacter := xp.Header[len(xp.Header)-1]
	if lastCharacter == ',' || lastCharacter == ';' {
		return XFCC{}, ex.New(ErrXFCCParsing).WithMessage("Ends with separator character")
	}

	for xp.Index < len(xp.Header) {
		char := xp.Header[xp.Index]
		switch xp.State {
		case parseXFCCKey:
			xp.HandleKeyCharacter(char)
		case parseXFCCValueStart:
			xp.HandleValueStartCharacter(char)
		case parseXFCCValue:
			err := xp.HandleValueCharacter(char)
			if err != nil {
				return XFCC{}, err
			}
		case parseXFCCValueQuoted:
			err := xp.HandleQuotedValueCharacter(char)
			if err != nil {
				return XFCC{}, err
			}
		}

		// Increment the index for the next iteration. (Note that branches of the
		// `switch` statement may have already incremented `i` as well.)
		xp.Index++
	}

	err := xp.Finalize()
	if err != nil {
		return XFCC{}, err
	}

	return xp.Parsed, nil
}

// FillXFCCKeyValue takes the currently active `.Key` and `.Value` and populates
// the current `.Element`.
func (xp *xfccParser) FillXFCCKeyValue() error {
	keyLower := strings.ToLower(string(xp.Key))
	switch keyLower {
	case "by":
		if xp.Element.By != "" {
			return ex.New(ErrXFCCParsing).WithMessagef("Key already encountered %q", keyLower)
		}
		xp.Element.By = string(xp.Value)
	case "hash":
		if len(xp.Element.Hash) > 0 {
			return ex.New(ErrXFCCParsing).WithMessagef("Key already encountered %q", keyLower)
		}
		xp.Element.Hash = string(xp.Value)
	case "cert":
		if len(xp.Element.Cert) > 0 {
			return ex.New(ErrXFCCParsing).WithMessagef("Key already encountered %q", keyLower)
		}
		xp.Element.Cert = string(xp.Value)
	case "chain":
		if len(xp.Element.Chain) > 0 {
			return ex.New(ErrXFCCParsing).WithMessagef("Key already encountered %q", keyLower)
		}
		xp.Element.Chain = string(xp.Value)
	case "subject":
		if len(xp.Element.Subject) > 0 {
			return ex.New(ErrXFCCParsing).WithMessagef("Key already encountered %q", keyLower)
		}
		xp.Element.Subject = string(xp.Value)
	case "uri":
		if xp.Element.URI != "" {
			return ex.New(ErrXFCCParsing).WithMessagef("Key already encountered %q", keyLower)
		}
		xp.Element.URI = string(xp.Value)
	case "dns":
		xp.Element.DNS = append(xp.Element.DNS, string(xp.Value))
	default:
		return ex.New(ErrXFCCParsing).WithMessagef("Unknown key %q", keyLower)
	}

	return nil
}

// HandleKeyCharacter advances the state machine if the current state is
// `parseXFCCKey`.
func (xp *xfccParser) HandleKeyCharacter(char rune) {
	if char == '=' {
		xp.State = parseXFCCValueStart
	} else {
		xp.Key = append(xp.Key, char)
	}
}

// HandleValueStartCharacter advances the state machine if the current state is
// `parseXFCCValueStart`.
func (xp *xfccParser) HandleValueStartCharacter(char rune) {
	if char == '"' {
		xp.State = parseXFCCValueQuoted
	} else {
		xp.Value = append(xp.Value, char)
		xp.State = parseXFCCValue
	}
}

// HandleValueCharacter advances the state machine if the current state is
// `parseXFCCValue`.
func (xp *xfccParser) HandleValueCharacter(char rune) error {
	if char == ',' || char == ';' {
		if len(xp.Key) == 0 || len(xp.Value) == 0 {
			return ex.New(ErrXFCCParsing).WithMessage("Key or Value missing")
		}
		err := xp.FillXFCCKeyValue()
		if err != nil {
			return err
		}

		xp.Key = make([]rune, 0, initialKeyCapacity)
		xp.Value = make([]rune, 0, initialValueCapacity)
		xp.State = parseXFCCKey
		if char == ',' {
			xp.Parsed = append(xp.Parsed, xp.Element)
			xp.Element = XFCCElement{}
		}
	} else {
		xp.Value = append(xp.Value, char)
	}

	return nil
}

// HandleQuotedValueCharacter advances the state machine if the current state is
// `parseXFCCValueQuoted`.
func (xp *xfccParser) HandleQuotedValueCharacter(char rune) error {
	if char == '\\' {
		nextIndex := xp.Index + 1
		if nextIndex < len(xp.Header) && xp.Header[nextIndex] == '"' {
			// Consume two characters at once here (since we have an
			// escaped quote).
			xp.Value = append(xp.Value, '"')
			xp.Index = nextIndex
		} else {
			xp.Value = append(xp.Value, char)
		}
	} else if char == '"' {
		// Since the **escaped quote** case `\"` has already been
		// covered, this case should only occur in the closing quote
		// case.
		nextIndex := xp.Index + 1
		if nextIndex < len(xp.Header) {
			if xp.Header[nextIndex] == ';' || xp.Header[nextIndex] == ',' {
				// Consume two characters at once here (since we have an
				// closing quote).
				xp.Index = nextIndex

				if len(xp.Key) == 0 {
					// Quoted values, e.g. `""`, are allowed to be empty.
					return ex.New(ErrXFCCParsing).WithMessage("Key missing")
				}
				err := xp.FillXFCCKeyValue()
				if err != nil {
					return err
				}

				xp.Key = make([]rune, 0, initialKeyCapacity)
				xp.Value = make([]rune, 0, initialValueCapacity)
				xp.State = parseXFCCKey
				if xp.Header[nextIndex] == ',' {
					xp.Parsed = append(xp.Parsed, xp.Element)
					xp.Element = XFCCElement{}
				}
			} else {
				return ex.New(ErrXFCCParsing).WithMessage("Closing quote not followed by `;`.")
			}
		} else {
			// NOTE: If `nextIndex >= len(xp.Header)` then we are at the end,
			//       which is a no-op here.
			xp.State = parseXFCCKey
		}
	} else {
		xp.Value = append(xp.Value, char)
	}

	return nil
}

// Finalize runs when the state machine has exhausted the `.Header`. It consumes
// any remaining `.Key` or `.Value` slices and adds them to the return value
// if need be.
func (xp *xfccParser) Finalize() error {
	if len(xp.Key) > 0 && len(xp.Value) > 0 {
		err := xp.FillXFCCKeyValue()
		if err != nil {
			return err
		}
	} else if len(xp.Key) > 0 || len(xp.Value) > 0 {
		return ex.New(ErrXFCCParsing).WithMessage("Key or value found but not both")
	}

	xp.Parsed = append(xp.Parsed, xp.Element)
	return nil
}
