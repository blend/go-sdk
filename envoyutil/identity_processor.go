/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package envoyutil

import (
	"fmt"
	"strings"

	"github.com/blend/go-sdk/collections"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/spiffeutil"
)

// NOTE: Ensure that
//   - `IdentityProcessor.KubernetesIdentityFormatter` satisfies `IdentityFormatter`
//   - `IdentityProcessor.IdentityProvider` satisfies `IdentityProvider`
var (
	_	IdentityFormatter	= IdentityProcessor{}.KubernetesIdentityFormatter
	_	IdentityProvider	= IdentityProcessor{}.IdentityProvider
)

// IdentityProcessorOption mutates an identity processor.
type IdentityProcessorOption func(*IdentityProcessor)

// OptIdentityType sets the identity type for the processor.
func OptIdentityType(it IdentityType) IdentityProcessorOption {
	return func(ip *IdentityProcessor) {
		ip.Type = it
	}
}

// OptAllowedTrustDomains adds allowed trust domains to the processor.
func OptAllowedTrustDomains(trustDomains ...string) IdentityProcessorOption {
	return func(ip *IdentityProcessor) {
		ip.AllowedTrustDomains = append(ip.AllowedTrustDomains, trustDomains...)
	}
}

// OptDeniedTrustDomains adds denied trust domains to the processor.
func OptDeniedTrustDomains(trustDomains ...string) IdentityProcessorOption {
	return func(ip *IdentityProcessor) {
		ip.DeniedTrustDomains = append(ip.DeniedTrustDomains, trustDomains...)
	}
}

// OptAllowedIdentities adds allowed identities to the processor.
func OptAllowedIdentities(identities ...string) IdentityProcessorOption {
	return func(ip *IdentityProcessor) {
		ip.AllowedIdentities = ip.AllowedIdentities.Union(
			collections.NewSet(identities...),
		)
	}
}

// OptDeniedIdentities adds denied identities to the processor.
func OptDeniedIdentities(identities ...string) IdentityProcessorOption {
	return func(ip *IdentityProcessor) {
		ip.DeniedIdentities = ip.DeniedIdentities.Union(
			collections.NewSet(identities...),
		)
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
	ClientIdentity	IdentityType	= 0
	// ServerIdentity represents server identity.
	ServerIdentity	IdentityType	= 1
)

// IdentityProcessor provides configurable fields that can be used to
// help validate a parsed SPIFFE URI and produce and validate an identity from
// a parsed SPIFFE URI. The `Type` field determines if a client or server
// identity should be provided; by default the type will be client identity.
type IdentityProcessor struct {
	Type			IdentityType
	AllowedTrustDomains	[]string
	DeniedTrustDomains	[]string
	AllowedIdentities	collections.Set[string]
	DeniedIdentities	collections.Set[string]
	FormatIdentity		IdentityFormatter
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
			Class:	ip.errInvalidIdentity(),
			XFCC:	xfcc.String(),
		}
	}

	pu, err := spiffeutil.Parse(uriValue)
	// NOTE: The `pu == nil` check is redundant, we expect `spiffeutil.Parse()`
	//       not to violate the invariant that `pu != nil` when `err == nil`.
	if err != nil || pu == nil {
		return "", &XFCCExtractionError{
			Class:	ip.errInvalidIdentity(),
			XFCC:	xfcc.String(),
		}
	}

	if err := ip.ProcessAllowedTrustDomains(xfcc, pu); err != nil {
		return "", err
	}
	if err := ip.ProcessDeniedTrustDomains(xfcc, pu); err != nil {
		return "", err
	}

	identity, err := ip.formatIdentity(xfcc, pu)
	if err != nil {
		return "", err
	}

	if err := ip.ProcessAllowedIdentities(xfcc, identity); err != nil {
		return "", err
	}
	if err := ip.ProcessDeniedIdentities(xfcc, identity); err != nil {
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
			Class:	ip.errInvalidIdentity(),
			XFCC:	xfcc.String(),
		}
	}
	return fmt.Sprintf("%s.%s", kw.ServiceAccount, kw.Namespace), nil
}

// ProcessAllowedTrustDomains returns an error if an allow list is configured
// and the trust domain from the parsed SPIFFE URI does not match any elements
// in the list.
func (ip IdentityProcessor) ProcessAllowedTrustDomains(xfcc XFCCElement, pu *spiffeutil.ParsedURI) error {
	if len(ip.AllowedTrustDomains) == 0 {
		return nil
	}

	for _, allowed := range ip.AllowedTrustDomains {
		if strings.EqualFold(pu.TrustDomain, allowed) {
			return nil
		}
	}
	return &XFCCValidationError{
		Class:	ip.errInvalidIdentity(),
		XFCC:	xfcc.String(),
		Metadata: map[string]string{
			"trustDomain": pu.TrustDomain,
		},
	}
}

// ProcessDeniedTrustDomains returns an error if a denied list is configured
// and the trust domain from the parsed SPIFFE URI matches any elements in the
// list.
func (ip IdentityProcessor) ProcessDeniedTrustDomains(xfcc XFCCElement, pu *spiffeutil.ParsedURI) error {
	for _, denied := range ip.DeniedTrustDomains {
		if strings.EqualFold(pu.TrustDomain, denied) {
			return &XFCCValidationError{
				Class:	ip.errInvalidIdentity(),
				XFCC:	xfcc.String(),
				Metadata: map[string]string{
					"trustDomain": pu.TrustDomain,
				},
			}
		}
	}

	return nil
}

// ProcessAllowedIdentities returns an error if an allow list is configured
// and the identity does not match any elements in the list.
func (ip IdentityProcessor) ProcessAllowedIdentities(xfcc XFCCElement, identity string) error {
	if ip.AllowedIdentities.Len() == 0 {
		return nil
	}

	if ip.AllowedIdentities.Contains(identity) {
		return nil
	}

	return &XFCCValidationError{
		Class:	ip.errDeniedIdentity(),
		XFCC:	xfcc.String(),
		Metadata: map[string]string{
			ip.getIdentityKey(): identity,
		},
	}
}

// ProcessDeniedIdentities returns an error if a denied list is configured
// and the identity matches any elements in the list.
func (ip IdentityProcessor) ProcessDeniedIdentities(xfcc XFCCElement, identity string) error {
	if ip.DeniedIdentities.Len() == 0 {
		return nil
	}

	if ip.DeniedIdentities.Contains(identity) {
		return &XFCCValidationError{
			Class:	ip.errDeniedIdentity(),
			XFCC:	xfcc.String(),
			Metadata: map[string]string{
				ip.getIdentityKey(): identity,
			},
		}
	}

	return nil
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

// getIdentityKey returns a key to be used in error metadata indicating if a
// value is client identity or server identity.
func (ip IdentityProcessor) getIdentityKey() string {
	if ip.Type == ClientIdentity {
		return "clientIdentity"
	}

	return "serverIdentity"
}

// errInvalidIdentity maps the `Type` to a specific error class indicating
// an invalid identity.
func (ip IdentityProcessor) errInvalidIdentity() ex.Class {
	if ip.Type == ClientIdentity {
		return ErrInvalidClientIdentity
	}

	return ErrInvalidServerIdentity
}

// errDeniedIdentity maps the `Type` to a specific error class indicating
// a denied identity.
func (ip IdentityProcessor) errDeniedIdentity() ex.Class {
	if ip.Type == ClientIdentity {
		return ErrDeniedClientIdentity
	}

	return ErrDeniedServerIdentity
}
