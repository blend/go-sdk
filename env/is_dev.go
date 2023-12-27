/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package env

// IsDev returns if the environment is development.
func IsDev(serviceEnv string) bool {
	switch serviceEnv {
	case ServiceEnvDev:
		return true
	default:
		return false
	}
}

// IsDevTest returns if the environment is a local development environment (i.e. `dev` or `test`).
func IsDevTest(serviceEnv string) bool {
	switch serviceEnv {
	case ServiceEnvDev, ServiceEnvTest:
		return true
	default:
		return false
	}
}

// IsDevlike returns if the environment is development.
// It is strictly the inverse of `IsProdlike`.
func IsDevlike(serviceEnv string) bool {
	return !IsProdlike(serviceEnv)
}
