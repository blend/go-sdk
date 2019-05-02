package certutil

import (
	"crypto/tls"
	"crypto/x509"
	"sync"

	"github.com/blend/go-sdk/ex"
)

// NewCertManagerWithKeyPairs returns a new cert pool from key pairs.
func NewCertManagerWithKeyPairs(server KeyPair, certificateAuthorities []KeyPair, clients ...KeyPair) (*CertManager, error) {
	serverCert, err := server.CertBytes()
	if err != nil {
		return nil, err
	}
	serverKey, err := server.KeyBytes()
	if err != nil {
		return nil, err
	}

	serverCertificate, err := tls.X509KeyPair(serverCert, serverKey)
	if err != nil {
		return nil, err
	}
	caCertPool, err := ExtendSystemCertPool(certificateAuthorities...)
	if err != nil {
		return nil, err
	}

	clientCerts := map[string][]byte{}
	for _, client := range clients {
		certPEM, err := client.CertBytes()
		if err != nil {
			return nil, err
		}
		commonNames, err := CommonNamesForCertPEM(certPEM)
		if err != nil {
			return nil, err
		}
		if len(commonNames) == 0 {
			return nil, ex.New(ErrInvalidCertPEM)
		}
		clientCerts[commonNames[0]] = certPEM
	}

	cm := NewCertManager(OptCertManagerServerCerts(serverCertificate), OptCertManagerRootCAs(caCertPool))
	return cm, cm.UpdateClientCerts(clientCerts)
}

// NewCertManager returns a new cert manager.
func NewCertManager(options ...CertManagerOption) *CertManager {
	certManager := &CertManager{
		TLSConfig: &tls.Config{
			ClientAuth: tls.RequireAndVerifyClientCert,
		},
		ClientCerts: map[string][]byte{},
	}
	certManager.TLSConfig.GetConfigForClient = certManager.GetConfigForClient

	for _, option := range options {
		option(certManager)
	}
	return certManager
}

// CertManagerOption is an option for a cert manager.
type CertManagerOption func(*CertManager)

// OptCertManagerRootCAs sets a field on the cert manager.
func OptCertManagerRootCAs(pool *x509.CertPool) CertManagerOption {
	return func(cm *CertManager) { cm.TLSConfig.RootCAs = pool }
}

// OptCertManagerServerCerts sets a field on the cert manager.
func OptCertManagerServerCerts(server ...tls.Certificate) CertManagerOption {
	return func(cm *CertManager) { cm.TLSConfig.Certificates = server }
}

// OptCertManagerClientCerts sets a field on the cert manager.
func OptCertManagerClientCerts(client *x509.CertPool) CertManagerOption {
	return func(cm *CertManager) { cm.TLSConfig.ClientCAs = client }
}

// CertManager is a pool of client certs.
type CertManager struct {
	sync.Mutex
	TLSConfig   *tls.Config
	ClientCerts map[string][]byte
}

// ClientCertUIDs returns all the client cert uids.
func (cm *CertManager) ClientCertUIDs() (output []string) {
	for uid := range cm.ClientCerts {
		output = append(output, uid)
	}
	return
}

// HasClientCert returns if the manager has a client cert.
func (cm *CertManager) HasClientCert(uid string) (has bool) {
	cm.Lock()
	_, has = cm.ClientCerts[uid]
	cm.Unlock()
	return
}

// AddClientCert adds a client cert to the bunde and refreshes the bundle.
func (cm *CertManager) AddClientCert(clientCert []byte) error {
	cm.Lock()
	defer cm.Unlock()

	commonNames, err := ParseCertPEM(clientCert)
	if err != nil {
		return err
	}
	if len(commonNames) == 0 {
		return ex.New(ErrInvalidCertPEM)
	}
	cm.ClientCerts[commonNames[0].Subject.CommonName] = clientCert
	return cm.RefreshClientCerts()
}

// RemoveClientCert removes a client cert by uid.
func (cm *CertManager) RemoveClientCert(uid string) error {
	cm.Lock()
	defer cm.Unlock()
	delete(cm.ClientCerts, uid)
	return cm.RefreshClientCerts()
}

// UpdateClientCerts sets the client cert bundle fully.
func (cm *CertManager) UpdateClientCerts(clientCerts map[string][]byte) error {
	cm.Lock()
	defer cm.Unlock()
	cm.ClientCerts = clientCerts
	return cm.RefreshClientCerts()
}

// RefreshClientCerts reloads the client cert bundle.
func (cm *CertManager) RefreshClientCerts() error {
	pool := x509.NewCertPool()
	for uid, cert := range cm.ClientCerts {
		if ok := pool.AppendCertsFromPEM(cert); !ok {
			return ex.New("invalid ca cert for client cert pool", ex.OptMessagef("cert uid: %s", uid))
		}
	}
	cm.TLSConfig.ClientCAs = pool
	cm.TLSConfig.BuildNameToCertificate()
	return nil
}

// GetConfigForClient gets a tls config for a given client hello.
func (cm *CertManager) GetConfigForClient(sni *tls.ClientHelloInfo) (config *tls.Config, _ error) {
	cm.Lock()
	config = cm.TLSConfig.Clone()
	cm.Unlock()
	return
}
