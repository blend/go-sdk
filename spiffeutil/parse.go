/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package spiffeutil

import (
	"strings"

	"github.com/blend/go-sdk/ex"
)

const (
	spiffePrefix	= "spiffe://"

	// ErrInvalidURI is the class of error returned when parsing SPIFFE URI fails
	ErrInvalidURI	= ex.Class("Invalid SPIFFE URI")
	// ErrNonKubernetesWorkload is class of error returned when parsing a SPIFFE
	// Kubernetes workload identifier fails.
	ErrNonKubernetesWorkload	= ex.Class("Workload ID not in Kubernetes format")
)

// ParsedURI represents a SPIFFE URI that has been parsed via `Parse()`.
type ParsedURI struct {
	TrustDomain	string
	WorkloadID	string
}

// KubernetesWorkload describes a Kubernetes workload identifier.
type KubernetesWorkload struct {
	Namespace	string
	ServiceAccount	string
}

// Parse consumes a SPIFFE URI and splits out the trust domain and workload
// identifier. For example in `spiffe://cluster.local/ns/blend/sa/quasar`
// the trust domain is `cluster.local` and the workload identifier is
// `ns/blend/sa/quasar`.
func Parse(uri string) (*ParsedURI, error) {
	if !strings.HasPrefix(uri, spiffePrefix) {
		return nil, ex.New(ErrInvalidURI).WithMessagef("Does not match protocol: %q", uri)
	}

	suffix := uri[len(spiffePrefix):]
	parts := strings.SplitN(suffix, "/", 2)
	if len(parts) != 2 || len(parts[1]) == 0 {
		return nil, ex.New(ErrInvalidURI).WithMessagef("Missing workload identifier: %q", uri)
	}

	pu := &ParsedURI{TrustDomain: parts[0], WorkloadID: parts[1]}
	return pu, nil
}

// ParseKubernetesWorkloadID parses a SPIFFE workload identifier that identifies
// Kubernetes service account, of the form  `ns/{namespace}/sa/{serviceAccount}`.
func ParseKubernetesWorkloadID(workloadID string) (*KubernetesWorkload, error) {
	parts := strings.Split(workloadID, "/")
	if len(parts) != 4 || parts[0] != "ns" || parts[2] != "sa" {
		err := ex.New(ErrNonKubernetesWorkload, ex.OptMessagef("Workload identifier: %q", workloadID))
		return nil, err
	}

	return &KubernetesWorkload{Namespace: parts[1], ServiceAccount: parts[3]}, nil
}
