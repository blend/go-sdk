/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package env

// Service specific constants
const (
	// VarRegion is a common env var name.
	VarRegion	= "AWS_REGION"
	// VarServiceEnv is a common env var name.
	VarServiceEnv	= "SERVICE_ENV"
	// VarServiceName is a common env var name.
	VarServiceName	= "SERVICE_NAME"
	// VarProjectName is a common env var name.
	VarProjectName	= "PROJECT_NAME"
	// VarClusterName is a common env var name.
	VarClusterName	= "CLUSTER_NAME"
	// VarServiceSecret is a common env var name.
	VarServiceSecret	= "SERVICE_SECRET"
	// VarPort is a common env var name.
	VarPort	= "PORT"
	// VarHostname is a common env var name.
	VarHostname	= "HOSTNAME"
	// VarVersion is a common env var name.
	VarVersion	= "VERSION"
	// VarGitRef is a common env var name.
	VarGitRef	= "GIT_REF"
	// VarSecurePort is a common env var name.
	VarSecurePort	= "SECURE_PORT"
	// VarTLSCertPath is a common env var name.
	VarTLSCertPath	= "TLS_CERT_PATH"
	// VarTLSKeyPath is a common env var name.
	VarTLSKeyPath	= "TLS_KEY_PATH"
	// VarTLSCert is a common env var name.
	VarTLSCert	= "TLS_CERT"
	// VarTLSKey is a common env var name.
	VarTLSKey	= "TLS_KEY"

	// VarPGIdleConns is a common env var name.
	VarPGIdleConns	= "PG_IDLE_CONNS"
	// VarPGMaxConns is a common env var name.
	VarPGMaxConns	= "PG_MAX_CONNS"

	// ServiceEnvTest is a service environment.
	ServiceEnvTest	= "test"
	// ServiceEnvSandbox is a service environment.
	ServiceEnvSandbox	= "sandbox"
	// ServiceEnvDev is a service environment.
	ServiceEnvDev	= "dev"
	// ServiceEnvCI is a service environment.
	ServiceEnvCI	= "ci"
	// ServiceEnvPreprod is a service environment.
	ServiceEnvPreprod	= "preprod"
	// ServiceEnvBeta is a service environment.
	ServiceEnvBeta	= "beta"
	// ServiceEnvProd is a service environment.
	ServiceEnvProd	= "prod"

	// DefaultServiceEnv is the default service env to use for configs.
	DefaultServiceEnv	= ServiceEnvDev

	// TagName is the reflection tag name.
	ReflectTagName	= "env"
)
