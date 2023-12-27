/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package env

// Unmarshaler is a type that implements `UnmarshalEnv`.
type Unmarshaler interface {
	UnmarshalEnv(vars Vars) error
}
