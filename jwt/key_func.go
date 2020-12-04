package jwt

// Keyfunc should return the key used in verification based on the raw token passed to it.
type Keyfunc func(*Token) (interface{}, error)

// KeyfuncStatic returns a static key func.
func KeyfuncStatic(key []byte) Keyfunc {
	return func(_ *Token) (interface{}, error) {
		return key, nil
	}
}
