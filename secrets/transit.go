package secrets

// Transit is an interface for an encryption-as-a-service client
type Transit interface {
	TransitEncrypt(key string, context map[string]interface{}, data []byte) (string, error)
	TransitDecrypt(key string, context map[string]interface{}, ciphertext string) ([]byte, error)
}
