package envoyutil

import (
	"fmt"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/spiffeutil"
)

// NOTE: Ensure that
//       - `IdentityProcessor.KubernetesIdentityFormatter` satisfies `IdentityFormatter`
//       - `IdentityProcessor.IdentityProvider` satisfies `IdentityProvider`
var (
	_ IdentityFormatter = IdentityProcessor{}.KubernetesIdentityFormatter
	_ IdentityProvider  = IdentityProcessor{}.IdentityProvider
)

// IdentityProcessorOption mutates an identity processor.
type IdentityProcessorOption func(*IdentityProcessor)

// OptIdentityType sets the identity type for the processor.
func OptIdentityType(it IdentityType) IdentityProcessorOption {
	return func(ip *IdentityProcessor) {
		ip.Type = it
	}
}

// OptFormatIdentity sets the `FormatIdentity` on the processor.
func OptFormatIdentity(formatter IdentityFormatter) IdentityProcessorOption {
	return func(ip *IdentityProcessor) {
		ip.FormatIdentity = formatter
	}
}

// IdentityFormatter describes functions that will produce an identity string
// from a parsed SPIFFE URI.
type IdentityFormatter = func(XFCCElement, *spiffeutil.ParsedURI) (string, error)

// IdentityType represents the type of identity that will be extracted by an
// `IdentityProcessor`. It can either be a client or server identity.
type IdentityType int

const (
	// ClientIdentity represents client identity.
	ClientIdentity IdentityType = 0
	// ServerIdentity represents server identity.
	ServerIdentity IdentityType = 1
)

// IdentityProcessor provides configurable fields that can be used to
// help validate a parsed SPIFFE URI and produce and validate an identity from
// a parsed SPIFFE URI. The `Type` field determines if a client or server
// identity should be provided; by default the type will be client identity.
type IdentityProcessor struct {
	Type           IdentityType
	FormatIdentity IdentityFormatter
}

// IdentityProvider returns a client or server identity; it uses the configured
// rules to validate and format the identity by parsing the `URI` field (for
// client identity) or `By` field (for server identity) of the XFCC element. If
// `FormatIdentity` has not been specified, the `KubernetesIdentityFormatter()`
// method will be used as a fallback.
//
// This method satisfies the `IdentityProvider` interface.
func (ip IdentityProcessor) IdentityProvider(xfcc XFCCElement) (string, error) {
	uriValue := ip.getURIForIdentity(xfcc)

	if uriValue == "" {
		return "", &XFCCValidationError{
			Class: ip.errInvalidIdentity(),
			XFCC:  xfcc.String(),
		}
	}

	pu, err := spiffeutil.Parse(uriValue)
	// NOTE: The `pu == nil` check is redundant, we expect `spiffeutil.Parse()`
	//       not to violate the invariant that `pu != nil` when `err == nil`.
	if err != nil || pu == nil {
		return "", &XFCCExtractionError{
			Class: ip.errInvalidIdentity(),
			XFCC:  xfcc.String(),
		}
	}

	identity, err := ip.formatIdentity(xfcc, pu)
	if err != nil {
		return "", err
	}

	return identity, nil
}

// KubernetesIdentityFormatter assumes the SPIFFE URI contains a Kubernetes
// workload ID of the form `ns/{namespace}/sa/{serviceAccount}` and formats the
// identity as `{serviceAccount}.{namespace}`. This function satisfies the
// `IdentityFormatter` interface.
func (ip IdentityProcessor) KubernetesIdentityFormatter(xfcc XFCCElement, pu *spiffeutil.ParsedURI) (string, error) {
	kw, err := spiffeutil.ParseKubernetesWorkloadID(pu.WorkloadID)
	if err != nil {
		return "", &XFCCExtractionError{
			Class: ip.errInvalidIdentity(),
			XFCC:  xfcc.String(),
		}
	}
	return fmt.Sprintf("%s.%s", kw.ServiceAccount, kw.Namespace), nil
}

// formatIdentity invokes the `FormatIdentity` on the current processor
// or falls back to `KubernetesIdentityFormatter()` if it is not set.
func (ip IdentityProcessor) formatIdentity(xfcc XFCCElement, pu *spiffeutil.ParsedURI) (string, error) {
	if ip.FormatIdentity != nil {
		return ip.FormatIdentity(xfcc, pu)
	}
	return ip.KubernetesIdentityFormatter(xfcc, pu)
}

// getURIForIdentity returns either the `URI` field if this processor has `Type`
// "client identity" or the `By` field for the server identity.
func (ip IdentityProcessor) getURIForIdentity(xfcc XFCCElement) string {
	if ip.Type == ClientIdentity {
		return xfcc.URI
	}

	return xfcc.By
}

// errInvalidIdentity maps the `Type` to a specific error class indicating
// an invalid identity.
func (ip IdentityProcessor) errInvalidIdentity() ex.Class {
	if ip.Type == ClientIdentity {
		return ErrInvalidClientIdentity
	}

	return ErrInvalidServerIdentity
}
