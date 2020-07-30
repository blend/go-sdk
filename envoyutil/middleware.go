package envoyutil

import (
	"github.com/blend/go-sdk/web"
)

const (
	// StateKeyClientIdentity is the key into a `web.Ctx` state holding the
	// client identity of the client calling through Envoy.
	StateKeyClientIdentity = "envoy-client-identity"
)

// GetClientIdentity returns the client identity of the calling service or
// `""` if the client identity is unset.
func GetClientIdentity(ctx *web.Ctx) string {
	typed, ok := ctx.StateValue(StateKeyClientIdentity).(string)
	if !ok {
		return ""
	}
	return typed
}
