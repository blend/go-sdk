package env

// ConfigBase is a base config you can use in your own config objects.
type ConfigBase struct {
	Name        string `yaml:"serviceName" env:"SERVICE_NAME"`
	Environment string `yaml:"serviceEnv,omitempty" env:"SERVICE_ENV"`
}

// Resolve adds extra resolution steps for the config.
func (c *ConfigBase) Resolve() error {
	return Env().ReadInto(c)
}

// NameOrDefault gets the service name.
func (c ConfigBase) NameOrDefault() string {
	if c.Name != "" {
		return c.Name
	}
	return ""
}

// EnvironmentOrDefault returns the service environment or a default.
func (c ConfigBase) EnvironmentOrDefault() string {
	if c.Environment != "" {
		return c.Environment
	}
	return DefaultServiceEnv
}

// IsProdlike returns if the cluster meta environment is prodlike.
func (c ConfigBase) IsProdlike() bool {
	env := c.EnvironmentOrDefault()
	return env != ServiceEnvDev &&
		env != ServiceEnvSandbox &&
		env != ServiceEnvTest &&
		env != ServiceEnvCI
}
