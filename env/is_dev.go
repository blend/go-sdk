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

// IsDevlike returns if the environment is development.
// It is strictly the inverse of `IsProdlike`.
func IsDevlike(serviceEnv string) bool {
	return !IsProdlike(serviceEnv)
}
