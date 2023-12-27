/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package oauth

var (
	// DefaultScopes is the default oauth scopes.
	DefaultScopes = []string{
		"openid",
		"email",
		"profile",
	}
)

const (
	// GoogleKeysURL is the url we fetch google's public verification keys in JWK form.
	GoogleKeysURL	= "https://www.googleapis.com/oauth2/v3/certs"
	// GoogleIssuer is the expected `iss` field on JWTs from google.
	GoogleIssuer	= "https://accounts.google.com"
	// GoogleIssuerAlternate is the alternate expected `iss` field on JWTs from google.
	GoogleIssuerAlternate	= "accounts.google.com"
)

const (
	// ErrCodeMissing is returned if the code was missing from an oauth return request.
	ErrCodeMissing	Error	= "state missing from request"
	// ErrStateMissing is returned if the state was missing from an oauth return request.
	ErrStateMissing	Error	= "state missing from request"
	// ErrInvalidHostedDomain is an error returned if the JWT hosted zone doesn't match any of the whitelisted domains.
	ErrInvalidHostedDomain	Error	= "hosted domain validation failed"
	// ErrInvalidAntiforgeryToken is an error returns on oauth finish that indicates we didn't originate the auth request.
	ErrInvalidAntiforgeryToken	Error	= "invalid anti-forgery token"

	// ErrInvalidJWTAudience is an error in validing the token jwt.
	ErrInvalidJWTAudience	Error	= "invalid jwt audience; should match clientID"
	// ErrInvalidJWTIssuer is an error in validing the token jwt.
	ErrInvalidJWTIssuer	Error	= "invalid jwt issuer; should be a valid google issuer"
	// ErrInvalidJWTHostedDomain is an error in validing the token jwt.
	ErrInvalidJWTHostedDomain	Error	= "invalid jwt hosted domain; must be in the allowed domain list"
	// ErrInvalidJWT is returned when we fail to decode or verify the token jwt.
	ErrInvalidJWT	Error	= "invalid jwt; failed to decode or verify"

	// ErrProfileJSONUnmarshal is an error returned if the json unmarshal failed.
	ErrProfileJSONUnmarshal	Error	= "profile json unmarshal failed"

	// ErrFailedCodeExchange happens if the code exchange for an access token fails.
	ErrFailedCodeExchange	Error	= "oauth code exchange failed"
	// ErrGoogleResponseStatus is an error that can occur when querying the google apis.
	ErrGoogleResponseStatus	Error	= "google returned a non 2xx response"

	// ErrSecretRequired is a configuration error indicating we did not provide a secret.
	ErrSecretRequired	Error	= "manager secret required"
	// ErrClientIDRequired is a self validation error.
	ErrClientIDRequired	Error	= "clientID is required"
	// ErrClientSecretRequired is a self validation error.
	ErrClientSecretRequired	Error	= "clientSecret is required"
	// ErrRedirectURIRequired is a self validation error.
	ErrRedirectURIRequired	Error	= "redirectURI is required"
	// ErrInvalidRedirectURI is an error in validating the redirect uri.
	ErrInvalidRedirectURI	Error	= "invalid redirectURI"
)
