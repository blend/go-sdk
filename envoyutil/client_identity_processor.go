package envoyutil

import (
	"fmt"
	"strings"

	"github.com/blend/go-sdk/collections"
	"github.com/blend/go-sdk/spiffeutil"
)

// NOTE: Ensure that
//       - `KubernetesClientIdentityFormatter` satisfies `ClientIdentityFormatter`
//       - `ClientIdentityProcessor.ClientIdentityProvider` satisfies `ClientIdentityProvider`
var (
	_ ClientIdentityFormatter = KubernetesClientIdentityFormatter
	_ ClientIdentityProvider  = ClientIdentityProcessor{}.ClientIdentityProvider
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
func OptAllowedClientIdentities(clientIdentities ...string) ClientIdentityProcessorOption {
	return func(cip *ClientIdentityProcessor) {
		cip.AllowedClientIdentities = cip.AllowedClientIdentities.Union(
			collections.NewSetOfString(clientIdentities...),
		)
	}
}

// OptDeniedClientIdentities adds denied client identities to the processor.
func OptDeniedClientIdentities(clientIdentities ...string) ClientIdentityProcessorOption {
	return func(cip *ClientIdentityProcessor) {
		cip.DeniedClientIdentities = cip.DeniedClientIdentities.Union(
			collections.NewSetOfString(clientIdentities...),
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

// ClientIdentityProvider returns a client identity; it uses the configured rules
// to validate and format the client identity by parsing the `URI` field of the
// XFCC element. If `FormatClientIdentity` has not been specified, the
// `KubernetesClientIdentityFormatter()` function will be used as a fallback.
//
// This function satisfies the `ClientIdentityProvider` interface.
func (cip ClientIdentityProcessor) ClientIdentityProvider(xfcc XFCCElement) (string, error) {
	pu, err := spiffeutil.Parse(xfcc.URI)
	// NOTE: The `pu == nil` check is redundant, we expect `spiffeutil.Parse()`
	//       not to violate the invariant that `pu != nil` when `err == nil`.
	if err != nil || pu == nil {
		return "", &XFCCExtractionError{
			Class: ErrInvalidClientIdentity,
			XFCC:  xfcc.String(),
		}
	}

	if err := cip.ProcessAllowedTrustDomains(xfcc, pu); err != nil {
		return "", err
	}
	if err := cip.ProcessDeniedTrustDomains(xfcc, pu); err != nil {
		return "", err
	}

	clientIdentity, err := cip.formatClientIdentity(xfcc, pu)
	if err != nil {
		return "", err
	}

	if err := cip.ProcessAllowedClientIdentities(xfcc, clientIdentity); err != nil {
		return "", err
	}
	if err := cip.ProcessDeniedClientIdentities(xfcc, clientIdentity); err != nil {
		return "", err
	}
	return clientIdentity, nil
}

// ProcessAllowedTrustDomains returns an error if an allow list is configured
// and the trust domain from the parsed SPIFFE URI does not match any elements
// in the list.
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
// and the trust domain from the parsed SPIFFE URI matches any elements in the
// list.
func (cip ClientIdentityProcessor) ProcessDeniedTrustDomains(xfcc XFCCElement, pu *spiffeutil.ParsedURI) error {
	for _, denied := range cip.DeniedTrustDomains {
		if strings.EqualFold(pu.TrustDomain, denied) {
			return &XFCCValidationError{
				Class: ErrInvalidClientIdentity,
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
// and the client identity does not match any elements in the list.
func (cip ClientIdentityProcessor) ProcessAllowedClientIdentities(xfcc XFCCElement, clientIdentity string) error {
	if cip.AllowedClientIdentities.Len() == 0 {
		return nil
	}

	if cip.AllowedClientIdentities.Contains(clientIdentity) {
		return nil
	}

	return &XFCCValidationError{
		Class: ErrDeniedClientIdentity,
		XFCC:  xfcc.String(),
		Metadata: map[string]string{
			"clientIdentity": clientIdentity,
		},
	}
}

// ProcessDeniedClientIdentities returns an error if a denied list is configured
// and the client identity matches any elements in the list.
func (cip ClientIdentityProcessor) ProcessDeniedClientIdentities(xfcc XFCCElement, clientIdentity string) error {
	if cip.DeniedClientIdentities.Len() == 0 {
		return nil
	}

	if cip.DeniedClientIdentities.Contains(clientIdentity) {
		return &XFCCValidationError{
			Class: ErrDeniedClientIdentity,
			XFCC:  xfcc.String(),
			Metadata: map[string]string{
				"clientIdentity": clientIdentity,
			},
		}
	}

	return nil
}

// formatClientIdentity invokes the `FormatClientIdentity` on the current processor
// or falls back to `KubernetesClientIdentityFormatter()` if it is not set.
func (cip ClientIdentityProcessor) formatClientIdentity(xfcc XFCCElement, pu *spiffeutil.ParsedURI) (string, error) {
	if cip.FormatClientIdentity != nil {
		return cip.FormatClientIdentity(xfcc, pu)
	}
	return KubernetesClientIdentityFormatter(xfcc, pu)
}
