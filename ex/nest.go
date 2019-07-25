package ex

// Nest nests an arbitrary number of exceptions.
func Nest(err ...error) error {
	var ex *Ex
	var last *Ex
	var didSet bool

	for _, e := range err {
		if e != nil {
			var wrappedEx *Ex
			if typedEx, isTyped := e.(*Ex); !isTyped {
				wrappedEx = &Ex{
					Class: e,
					Stack: Callers(DefaultStartDepth),
				}
			} else {
				wrappedEx = typedEx
			}

			if wrappedEx != ex {
				if ex == nil {
					ex = wrappedEx
					last = wrappedEx
				} else {
					last.Inner = wrappedEx
					last = wrappedEx
				}
				didSet = true
			}
		}
	}
	if didSet {
		return ex
	}
	return nil
}
