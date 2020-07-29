package envoyutil

import (
	"fmt"
	"strings"

	"github.com/blend/go-sdk/collections"
	"github.com/blend/go-sdk/spiffeutil"
)

// NOTE: Ensure that
//       - `KubernetesClientIdentityFormatter` satisfies `ClientIdentityFormatter`
var (
	_ ClientIdentityFormatter = KubernetesClientIdentityFormatter
)

// OptAllowedTrustDomains adds allowed trust domains to the processor.
func OptAllowedTrustDomains(trustDomains ...string) ClientIdentityProcessorOption {
	return func(cip *ClientIdentityProcessor) {
		cip.AllowedTrustDomains = append(cip.AllowedTrustDomains, trustDomains...)
	}
}

// OptDeniedTrustDomains adds denied trust domains to the processor.
func OptDeniedTrustDomains(trustDomains ...string) ClientIdentityProcessorOption {
	return func(cip *ClientIdentityProcessor) {
		cip.DeniedTrustDomains = append(cip.DeniedTrustDomains, trustDomains...)
	}
}

// OptAllowedClientIdentities adds allowed client identities to the processor.
func OptAllowedClientIdentities(clientIDs ...string) ClientIdentityProcessorOption {
	return func(cip *ClientIdentityProcessor) {
		cip.AllowedClientIdentities = cip.AllowedClientIdentities.Union(
			collections.NewSetOfString(clientIDs...),
		)
	}
}

// OptDeniedClientIdentities adds denied client identities to the processor.
func OptDeniedClientIdentities(clientIDs ...string) ClientIdentityProcessorOption {
	return func(cip *ClientIdentityProcessor) {
		cip.DeniedClientIdentities = cip.DeniedClientIdentities.Union(
			collections.NewSetOfString(clientIDs...),
		)
	}
}

// OptFormatClientIdentity sets the `FormatClientIdentity` on the processor.
func OptFormatClientIdentity(formatter ClientIdentityFormatter) ClientIdentityProcessorOption {
	return func(cip *ClientIdentityProcessor) {
		cip.FormatClientIdentity = formatter
	}
}

// ClientIdentityProcessorOption mutates a client identity processor.
type ClientIdentityProcessorOption func(*ClientIdentityProcessor)

// KubernetesClientIdentityFormatter assumes the SPIFFE URI contains a Kubernetes
// workload ID of the form `ns/{namespace}/sa/{serviceAccount}` and formats the
// client identity as `{serviceAccount}.{namespace}`. This function satisfies the
// `ClientIdentityFormatter` interface.
func KubernetesClientIdentityFormatter(xfcc XFCCElement, pu *spiffeutil.ParsedURI) (string, error) {
	kw, err := spiffeutil.ParseKubernetesWorkloadID(pu.WorkloadID)
	if err != nil {
		return "", &XFCCExtractionError{
			Class: ErrInvalidClientIdentity,
			XFCC:  xfcc.String(),
		}
	}
	return fmt.Sprintf("%s.%s", kw.ServiceAccount, kw.Namespace), nil
}

// ClientIdentityFormatter describes functions that will produce a client
// identity string from a parsed SPIFFE URI.
type ClientIdentityFormatter = func(XFCCElement, *spiffeutil.ParsedURI) (string, error)

// ClientIdentityProcessor provides configurable fields that can be used to
// help validate a parsed SPIFFE URI and produce and validate a client identity
// from a parsed SPIFFE URI.
type ClientIdentityProcessor struct {
	AllowedTrustDomains     []string
	DeniedTrustDomains      []string
	AllowedClientIdentities collections.SetOfString
	DeniedClientIdentities  collections.SetOfString
	FormatClientIdentity    ClientIdentityFormatter
}

// ClientIdentityProvider returns a client identity provider for the given rule options.
// If a `WorkloadFormatter` has not been specified, the `KubernetesClientIdentityFormatter()`
// function will be used as a fallback.
func (cip ClientIdentityProcessor) ClientIdentityProvider(xfcc XFCCElement) (string, error) {
	pu, err := spiffeutil.Parse(xfcc.URI)
	if err != nil {
		return "", err
	}
	if err := cip.ProcessAllowedTrustDomains(xfcc, pu); err != nil {
		return "", err
	}
	if err := cip.ProcessDeniedTrustDomains(xfcc, pu); err != nil {
		return "", err
	}

	clientID, err := cip.formatClientID(xfcc, pu)
	if err != nil {
		return "", err
	}

	if err := cip.ProcessAllowedClientIdentities(xfcc, clientID); err != nil {
		return "", err
	}
	if err := cip.ProcessDeniedClientIdentities(xfcc, clientID); err != nil {
		return "", err
	}
	return clientID, nil
}

// ProcessAllowedTrustDomains returns an error if an allow list is configured
// and a trust domain does not match any elements in the list.
func (cip ClientIdentityProcessor) ProcessAllowedTrustDomains(xfcc XFCCElement, pu *spiffeutil.ParsedURI) error {
	if len(cip.AllowedTrustDomains) == 0 {
		return nil
	}

	for _, allowed := range cip.AllowedTrustDomains {
		if strings.EqualFold(pu.TrustDomain, allowed) {
			return nil
		}
	}
	return &XFCCValidationError{
		Class: ErrInvalidClientIdentity,
		XFCC:  xfcc.String(),
		Metadata: map[string]string{
			"trustDomain": pu.TrustDomain,
		},
	}
}

// ProcessDeniedTrustDomains returns an error if a denied list is configured
// and a trust domain matches any elements in the list.
func (cip ClientIdentityProcessor) ProcessDeniedTrustDomains(xfcc XFCCElement, pu *spiffeutil.ParsedURI) error {
	for _, denied := range cip.DeniedTrustDomains {
		if strings.EqualFold(pu.TrustDomain, denied) {
			return &XFCCValidationError{
				Class: ErrDeniedClientIdentity,
				XFCC:  xfcc.String(),
				Metadata: map[string]string{
					"trustDomain": pu.TrustDomain,
				},
			}
		}
	}

	return nil
}

// ProcessAllowedClientIdentities returns an error if an allow list is configured
// and a client ID does not match any elements in the list.
func (cip ClientIdentityProcessor) ProcessAllowedClientIdentities(xfcc XFCCElement, clientID string) error {
	if cip.AllowedClientIdentities.Len() == 0 {
		return nil
	}

	if cip.AllowedClientIdentities.Contains(clientID) {
		return nil
	}

	return &XFCCValidationError{
		Class: ErrInvalidClientIdentity,
		XFCC:  xfcc.String(),
		Metadata: map[string]string{
			"clientID": clientID,
		},
	}
}

// ProcessDeniedClientIdentities returns an error if a denied list is configured
// and a client ID matches any elements in the list.
func (cip ClientIdentityProcessor) ProcessDeniedClientIdentities(xfcc XFCCElement, clientID string) error {
	if cip.DeniedClientIdentities.Len() == 0 {
		return nil
	}

	if cip.DeniedClientIdentities.Contains(clientID) {
		return &XFCCValidationError{
			Class: ErrInvalidClientIdentity,
			XFCC:  xfcc.String(),
			Metadata: map[string]string{
				"clientID": clientID,
			},
		}
	}

	return nil
}

func (cip ClientIdentityProcessor) formatClientID(xfcc XFCCElement, pu *spiffeutil.ParsedURI) (string, error) {
	if cip.FormatClientIdentity != nil {
		return cip.FormatClientIdentity(xfcc, pu)
	}
	return KubernetesClientIdentityFormatter(xfcc, pu)
}
