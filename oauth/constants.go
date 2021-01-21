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
	// ErrCodeMissing is returned if the code was missing from an oauth return request.
	ErrCodeMissing Error = "state missing from request"
	// ErrStateMissing is returned if the state was missing from an oauth return request.
	ErrStateMissing Error = "state missing from request"
	// ErrInvalidHostedDomain is an error returned if the JWT hosted zone doesn't match any of the whitelisted domains.
	ErrInvalidHostedDomain Error = "hosted domain validation failed"
	// ErrInvalidAntiforgeryToken is an error returns on oauth finish that indicates we didn't originate the auth request.
	ErrInvalidAntiforgeryToken Error = "invalid anti-forgery token"

	// ErrFailedCodeExchange happens if the code exchange for an access token fails.
	ErrFailedCodeExchange Error = "oauth code exchange failed"
	// ErrGoogleResponseStatus is an error that can occur when querying the google apis.
	ErrGoogleResponseStatus Error = "google returned a non 2xx response"

	// ErrProfileJSONUnmarshal is an error returned if the json unmarshal failed.
	ErrProfileJSONUnmarshal Error = "profile json unmarshal failed"

	// ErrSecretRequired is a configuration error indicating we did not provide a secret.
	ErrSecretRequired Error = "manager secret required"
	// ErrClientIDRequired is a self validation error.
	ErrClientIDRequired Error = "clientID is required"
	// ErrClientSecretRequired is a self validation error.
	ErrClientSecretRequired Error = "clientSecret is required"
	// ErrRedirectURIRequired is a self validation error.
	ErrRedirectURIRequired Error = "redirectURI is required"
	// ErrInvalidRedirectURI is an error in validating the redirect uri.
	ErrInvalidRedirectURI Error = "invalid redirectURI"
)
