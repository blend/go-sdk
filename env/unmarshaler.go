package env

// Unmarshaler is a type that implements `UnmarshalEnv`.
type Unmarshaler interface {
	UnmarshalEnv(vars Vars) error
}
