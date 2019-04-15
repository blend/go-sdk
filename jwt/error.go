package jwt

import "github.com/blend/go-sdk/ex"

// Error constants.
var (
	// ErrValidation will be the top most class in most cases.
	ErrValidation ex.Class = "validation error"

	ErrValidationAudienceUnset ex.Class = "token claims audience unset"
	ErrValidationExpired       ex.Class = "token expired"
	ErrValidationIssued        ex.Class = "token issued in future"
	ErrValidationNotBefore     ex.Class = "token not before"

	ErrValidationSignature ex.Class = "signature is invalid"

	ErrKeyfuncUnset         ex.Class = "keyfunc is unset"
	ErrInvalidKey           ex.Class = "key is invalid"
	ErrInvalidKeyType       ex.Class = "key is of invalid type"
	ErrInvalidSigningMethod ex.Class = "invalid signing method"
	ErrHashUnavailable      ex.Class = "the requested hash function is unavailable"

	ErrHMACSignatureInvalid ex.Class = "hmac signature is invalid"

	ErrECDSAVerification ex.Class = "crypto/ecdsa: verification error"

	ErrKeyMustBePEMEncoded ex.Class = "invalid key: key must be pem encoded pkcs1 or pkcs8 private key"
	ErrNotRSAPrivateKey    ex.Class = "key is not a valid rsa private key"
	ErrNotRSAPublicKey     ex.Class = "key is not a valid rsa public key"
)

// IsValidation returns if the error is a validation error
// instead of a more structural error with the key infrastructure.
func IsValidation(err error) bool {
	return ex.Is(err, ErrValidation)
}
