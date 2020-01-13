package secrets

import (
	"context"

	"github.com/blend/go-sdk/ex"
)

const errNilConfig ex.Class = "config cannot be nil"

// SecretTraceConfig are the options for sending trace messages for the secrets package
type SecretTraceConfig struct {
	VaultOperation string
	KeyName        string
}

// TraceOption is an option type for secret trace
type TraceOption func(config *SecretTraceConfig) error

// OptTraceConfig allows you to provide the entire secret trace configuration
func OptTraceConfig(providedConfig SecretTraceConfig) TraceOption {
	return func(config *SecretTraceConfig) error {
		*config = providedConfig
		return nil
	}
}

// OptTraceVaultOperation allows you to set the VaultOperation being hit
func OptTraceVaultOperation(path string) TraceOption {
	return func(config *SecretTraceConfig) error {
		if config == nil {
			return errNilConfig
		}
		config.VaultOperation = path
		return nil
	}
}

// OptTraceKeyName allows you to specify the name of the key being interacted with
func OptTraceKeyName(keyName string) TraceOption {
	return func(config *SecretTraceConfig) error {
		if config == nil {
			return errNilConfig
		}
		config.KeyName = keyName
		return nil
	}
}

// Tracer is a tracer for requests.
type Tracer interface {
	Start(ctx context.Context, options ...TraceOption) (TraceFinisher, error)
}

// TraceFinisher is a finisher for traces.
type TraceFinisher interface {
	Finish(ctx context.Context, statusCode int, vaultError error)
}
