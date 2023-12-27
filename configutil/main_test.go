/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package configutil

import (
	"context"
	"os"
	"testing"

	"github.com/blend/go-sdk/env"
)

// TestMain is the testing entrypoint.
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

type config struct {
	Environment	string	`json:"env" yaml:"env" env:"SERVICE_ENV"`
	Other		string	`json:"other" yaml:"other" env:"OTHER"`
	Base		string	`json:"base" yaml:"base"`
}

type resolvedConfig struct {
	config
}

// Resolve implements configutil.BareResolver.
func (r *resolvedConfig) Resolve(ctx context.Context) error {
	r.Environment = env.GetVars(ctx).String("ENVIRONMENT")
	return nil
}

type fullConfig struct {
	fullConfigMeta	`yaml:",inline"`

	Field0	string	`json:"field0" yaml:"field0"`
	Field1	string	`json:"field1" yaml:"field1"`
	Field2	string	`json:"field2" yaml:"field2"`
	Field3	string	`json:"field3" yaml:"field3"`

	Child	fullConfigChild	`json:"child" yaml:"child"`
}

func (fc *fullConfig) Resolve(ctx context.Context) error {
	return Resolve(ctx,
		(&fc.fullConfigMeta).Resolve,
		(&fc.Child).Resolve,

		SetString(&fc.Field0, Env("CONFIGUTIL_FIELD0"), String(fc.Field0), String("default-field0")),
		SetString(&fc.Field1, Env("CONFIGUTIL_FIELD1"), String(fc.Field1), String("default-field1")),
		SetString(&fc.Field2, Env("CONFIGUTIL_FIELD2"), String(fc.Field2), String("default-field2")),
		SetString(&fc.Field3, Env("CONFIGUTIL_FIELD3"), String(fc.Field3), String("default-field3")),
	)
}

type fullConfigMeta struct {
	ServiceName	string	`json:"serviceName" yaml:"serviceName"`
	ServiceEnv	string	`json:"serviceEnv" yaml:"serviceEnv"`
	Version		string	`json:"version" yaml:"version"`
}

func (fcm *fullConfigMeta) Resolve(ctx context.Context) error {
	return Resolve(ctx,
		SetString(&fcm.ServiceEnv, Env("SERVICE_ENV"), String(fcm.ServiceEnv), String("dev")),
		SetString(&fcm.ServiceName, Env("SERVICE_NAME"), String(fcm.ServiceName), String("configutil")),
		SetString(&fcm.Version, Env("VERSION"), String(fcm.Version), String("0.1.0")),
	)
}

type fullConfigChild struct {
	Field0	string	`json:"field0" yaml:"field0"`
	Field1	string	`json:"field1" yaml:"field1"`
	Field2	string	`json:"field2" yaml:"field2"`
	Field3	string	`json:"field3" yaml:"field3"`
}

func (fcc *fullConfigChild) Resolve(ctx context.Context) error {
	return Resolve(ctx,
		SetString(&fcc.Field0, Env("CONFIGUTIL_CHILD_FIELD0"), String(fcc.Field0), String("default-child-field0")),
		SetString(&fcc.Field1, Env("CONFIGUTIL_CHILD_FIELD1"), String(fcc.Field1), String("default-child-field1")),
		SetString(&fcc.Field2, Env("CONFIGUTIL_CHILD_FIELD2"), String(fcc.Field2), String("default-child-field2")),
		SetString(&fcc.Field3, Env("CONFIGUTIL_CHILD_FIELD3"), String(fcc.Field3), String("default-child-field3")),
	)
}
