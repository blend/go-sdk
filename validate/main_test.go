package validate

func none() error { return nil }

func some(err error) func() error { return func() error { return err } }
