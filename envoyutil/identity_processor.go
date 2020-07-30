package envoyutil

import (
	"fmt"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/spiffeutil"
)

// NOTE: Ensure that
//       - `IdentityProcessor.KubernetesIdentityFormatter` satisfies `IdentityFormatter`
var (
	_ IdentityFormatter = IdentityProcessor{}.KubernetesIdentityFormatter
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

// errInvalidIdentity maps the `Type` to a specific error class indicating
// an invalid identity.
func (ip IdentityProcessor) errInvalidIdentity() ex.Class {
	if ip.Type == ClientIdentity {
		return ErrInvalidClientIdentity
	}

	return ErrInvalidServerIdentity
}
