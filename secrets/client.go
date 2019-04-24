package secrets

// Client is the general interface for a Secrets client
type Client interface {
	KVClient
	TransitClient
}
