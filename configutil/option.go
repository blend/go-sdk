package configutil

// Option is a modification of config options.
type Option func(*ConfigOptions) error

// OptAddPaths adds paths to search for the config file.
func OptAddPaths(paths ...string) Option {
	return func(co *ConfigOptions) error {
		co.Paths = append(co.Paths, paths...)
		return nil
	}
}

// OptAddPreferredPaths adds paths to search first for the config file.
func OptAddPreferredPaths(paths ...string) Option {
	return func(co *ConfigOptions) error {
		co.Paths = append(paths, co.Paths...)
		return nil
	}
}

// OptPaths sets paths to search for the config file.
func OptPaths(paths ...string) Option {
	return func(co *ConfigOptions) error {
		co.Paths = paths
		return nil
	}
}

// OptResolver sets an additional resolver for the config read.
func OptResolver(resolver func(interface{}) error) Option {
	return func(co *ConfigOptions) error {
		co.Resolver = resolver
		return nil
	}
}
