package testutil

// OptBefore appends before run actions.
func OptBefore(steps ...SuiteAction) Option {
	return func(s *Suite) {
		s.Before = append(s.Before, steps...)
	}
}
