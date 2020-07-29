package envoyutil

import (
	"fmt"
	"strings"

	"github.com/blend/go-sdk/spiffeutil"
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

// ClientIdentityProcessorOption mutates a client identity processor.
type ClientIdentityProcessorOption func(*ClientIdentityProcessor)

// DefaultWorkloadFormatter is a default workload formatter.
func DefaultWorkloadFormatter(xfcc XFCCElement, pu *spiffeutil.ParsedURI) (string, error) {
	kw, err := spiffeutil.ParseKubernetesWorkloadID(pu.WorkloadID)
	if err != nil {
		return "", &XFCCExtractionError{
			Class: ErrInvalidClientIdentity,
			XFCC:  xfcc.String(),
		}
	}
	return fmt.Sprintf("%s.%s", kw.ServiceAccount, kw.Namespace), nil
}

// ClientIdentityProcessor is a client identity processor.
type ClientIdentityProcessor struct {
	AllowedTrustDomains []string
	DeniedTrustDomains  []string
	WorkloadFormatter   func(XFCCElement, *spiffeutil.ParsedURI) (string, error)
}

// ClientIdentityProvider returns a client identity provider for the given rule options.
func (cip ClientIdentityProcessor) ClientIdentityProvider(xfcc XFCCElement) (string, error) {
	pu, err := spiffeutil.Parse(xfcc.URI)
	if err != nil {
		return "", err
	}
	if err := cip.ProcessAllowed(xfcc, pu); err != nil {
		return "", err
	}
	if err := cip.ProcessDenied(xfcc, pu); err != nil {
		return "", err
	}
	if cip.WorkloadFormatter != nil {
		return cip.WorkloadFormatter(xfcc, pu)
	}
	return DefaultWorkloadFormatter(xfcc, pu)
}

// ProcessAllowed returns an error if an allow list is configured and a trust domain does not match
// any elements in the list.
func (cip ClientIdentityProcessor) ProcessAllowed(xfcc XFCCElement, pu *spiffeutil.ParsedURI) error {
	if len(cip.AllowedTrustDomains) > 0 {
		for _, allowed := range cip.AllowedTrustDomains {
			if strings.EqualFold(pu.TrustDomain, allowed) {
				return nil
			}
		}
		return &XFCCValidationError{
			Class: ErrInvalidClientIdentity,
			XFCC:  xfcc.String(),
		}
	}
	return nil
}

// ProcessDenied returns an error if a denied list is configured and a trust domain matches
// any elements in the list.
func (cip ClientIdentityProcessor) ProcessDenied(xfcc XFCCElement, pu *spiffeutil.ParsedURI) error {
	if len(cip.DeniedTrustDomains) > 0 {
		for _, denied := range cip.DeniedTrustDomains {
			if strings.EqualFold(pu.TrustDomain, denied) {
				return &XFCCValidationError{
					Class: ErrDeniedClientIdentity,
					XFCC:  xfcc.String(),
				}
			}
		}
	}
	return nil
}
