/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package envoyutil

import (
	"context"
	"net/http"

	"github.com/blend/go-sdk/web"
)

// clientIdentityKey is an unexported sentinel type used to set and get the
// client identity string on a `context`.
type clientIdentityKey struct{}

// WithClientIdentity adds the client identity to a context.
//
// For example if `rc` is a `*web.Ctx`
// ```
// rc.WithContext(envoyutil.WithClientIdentity(rc.Context(), clientIdentity))
// ```
func WithClientIdentity(ctx context.Context, clientIdentity string) context.Context {
	return context.WithValue(ctx, clientIdentityKey{}, clientIdentity)
}

// GetClientIdentity returns the client identity of the calling service or
// `""` if the client identity is unset.
//
// For example if `rc` is a `*web.Ctx`
// ```
// envoyutil.GetClientIdentity(rc.Context())
// ```
func GetClientIdentity(ctx context.Context) string {
	value := ctx.Value(clientIdentityKey{})
	asStr, _ := value.(string)
	// If the type assertion fails, we still want to return the zero value for string.
	return asStr
}

// ClientIdentityRequired produces a middleware function that determines the
// client identity used in a connection secured with mTLS.
//
// This parses the `X-Forwarded-Client-Cert` (XFCC) from a request and uses
// a client identity provider (`cip`, e.g. see `SPIFFEClientIdentityProvider()`) to
// determine the client identity. Additionally, optional `verifiers` (e.g. see
// `SPIFFEServerIdentityProvider()`) can be used to verify other parts of the XFCC
// header such as the identity of the current server.
//
// In cases of error, the client identity will not be set on the current
// context. For error status codes 400 and 401, the error will be serialized as
// JSON or XML (via `ctx.DefaultProvider`) and returned in the HTTP response.
// For error status code 500, no identifying information from the error will be
// returned in the HTTP response.
//
// A 401 Unauthorized will be returned in the following cases:
//   - The XFCC header is missing
//   - The XFCC header (after parsing) contains zero elements or multiple elements
//     (this code expects exactly one XFCC element, under the assumption that the
//     Envoy `ForwardClientCertDetails` setting is configured to `SANITIZE_SET`)
//   - The values from the XFCC header fail custom validation provided by `cip` or
//     `verifiers`. For example, if the client identity is contained in a deny
//     list, this would be considered a validation error.
//
// A 400 Bad Request will be returned in the following cases:
//   - The XFCC header cannot be parsed
//   - Custom parsing / extraction done by `cip` fails. For example, in cases
//     where the `URI` field in the XFCC is expected to be a valid SPIFFE URI
//     with a valid Kubernetes workload identifier, if the `URI` field does
//     not follow that format (e.g. `urn:uuid:6e8bc430-9c3a-11d9-9669-0800200c9a66`)
//     this would be considered an extraction error.
//
// A 500 Internal Server Error will be returned if the error is unrelated to
// validating the XFCC header or to parsing / extracting values from the XFCC
// header.
func ClientIdentityRequired(cip IdentityProvider, verifiers ...VerifyXFCC) web.Middleware {
	return func(action web.Action) web.Action {
		return func(ctx *web.Ctx) web.Result {
			clientIdentity, err := ExtractAndVerifyClientIdentity(ctx.Request, cip, verifiers...)
			if IsValidationError(err) {
				return ctx.DefaultProvider.Status(http.StatusUnauthorized, err)
			}
			if IsExtractionError(err) {
				// NOTE: We don't use `ctx.DefaultProvider.BadRequest()` because
				//       we want to allow serializing `err` as JSON if possible.
				//       The JSON provider just uses `err.Error()` for the response.
				return ctx.DefaultProvider.Status(http.StatusBadRequest, err)
			}
			if err != nil {
				return ctx.DefaultProvider.InternalError(nil)
			}

			ctx.WithContext(WithClientIdentity(ctx.Context(), clientIdentity))
			return action(ctx)
		}
	}
}

// ClientIdentityAware produces a middleware function nearly identical to
// `ClientIdentityRequired`. The primary difference is that this middleware will
// **not** return an error HTTP response for extraction or validation errors;
// it will still return a 500 Internal Server Error in unexpected failures.
// In cases of extraction or validation errors, the middleware will pass along
// to the next `action` and the client identity is not set on the current context.
func ClientIdentityAware(cip IdentityProvider, verifiers ...VerifyXFCC) web.Middleware {
	return func(action web.Action) web.Action {
		return func(ctx *web.Ctx) web.Result {
			clientIdentity, err := ExtractAndVerifyClientIdentity(ctx.Request, cip, verifiers...)
			// Early exit for a no-op in cases of validation or extraction error.
			if IsValidationError(err) || IsExtractionError(err) {
				return action(ctx)
			}

			if err != nil {
				return ctx.DefaultProvider.InternalError(nil)
			}

			currentCtx := ctx.Context()
			ctx.WithContext(WithClientIdentity(currentCtx, clientIdentity))
			return action(ctx)
		}
	}
}
