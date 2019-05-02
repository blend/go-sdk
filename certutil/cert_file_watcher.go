package certutil

import (
	"crypto/tls"
	"os"
	"sync"
	"time"

	"github.com/blend/go-sdk/async"
	"github.com/blend/go-sdk/ex"
)

// Error constants.
const (
	ErrTLSPathsUnset ex.Class = "tls cert or key path unset; cannot continue"
)

// NewCertFileWatcher creates a new CertReloader object with a reload delay
func NewCertFileWatcher(certPath, keyPath string, opts ...CertFileWatcherOption) (*CertFileWatcher, error) {
	if certPath == "" || keyPath == "" {
		return nil, ex.New(ErrTLSPathsUnset)
	}

	cw := &CertFileWatcher{
		Latch:    async.NewLatch(),
		CertPath: certPath,
		KeyPath:  keyPath,
	}

	for _, opt := range opts {
		if err := opt(cw); err != nil {
			return nil, err
		}
	}

	// load cert to make sure the current key pair is valid
	if err := cw.Reload(); err != nil {
		return nil, err
	}
	return cw, nil
}

// CertFileWatcherOption is an option for a cert watcher.
type CertFileWatcherOption func(*CertFileWatcher) error

// CertFileWatcher reloads a cert key pair when there is a change, e.g. cert renewal
type CertFileWatcher struct {
	*async.Latch
	syncRoot sync.RWMutex

	Certificate *tls.Certificate

	CertPath     string
	KeyPath      string
	PollInterval time.Duration

	OnReload func()
}

// PollIntervalOrDefault returns the polling interval or a default.
func (cw *CertFileWatcher) PollIntervalOrDefault() time.Duration {
	if cw.PollInterval > 0 {
		return cw.PollInterval
	}
	return 500 * time.Millisecond
}

// Reload forces the reload of the underlying certificate.
func (cw *CertFileWatcher) Reload() error {
	cw.syncRoot.Lock()
	defer cw.syncRoot.Unlock()

	cert, err := tls.LoadX509KeyPair(cw.CertPath, cw.KeyPath)
	if err != nil {
		return ex.New(err)
	}
	cw.Certificate = &cert

	if cw.OnReload != nil {
		cw.OnReload()
	}
	return nil
}

// GetCertificate gets the cached certificate, it blocks when the `cert` field is being updated
func (cw *CertFileWatcher) GetCertificate(_ *tls.ClientHelloInfo) (*tls.Certificate, error) {
	cw.syncRoot.RLock()
	defer cw.syncRoot.RUnlock()
	return cw.Certificate, nil
}

// Start watches the cert and triggers a reload on change
func (cw *CertFileWatcher) Start() error {
	keyStat, err := os.Stat(cw.KeyPath)
	if err != nil {
		return err
	}
	certStat, err := os.Stat(cw.CertPath)
	if err != nil {
		return err
	}

	cw.Started()

	keyLastMod := keyStat.ModTime()
	certLastMod := certStat.ModTime()
	var didReload bool

	ticker := time.Tick(cw.PollIntervalOrDefault())
	for {
		select {
		case <-ticker:
			didReload = false

			// check key
			keyStat, err = os.Stat(cw.KeyPath)
			if err != nil {
				return err
			}
			if keyStat.ModTime().After(keyLastMod) {
				if err := cw.Reload(); err != nil {
					return err
				}
				didReload = true
				keyLastMod = keyStat.ModTime()
			}

			// check cert
			certStat, err = os.Stat(cw.CertPath)
			if err != nil {
				return err
			}
			if certStat.ModTime().After(certLastMod) {
				if !didReload {
					if err := cw.Reload(); err != nil {
						return err
					}
					didReload = true
				}
				certLastMod = certStat.ModTime()
			}
		case <-cw.NotifyStopping():
			cw.Stopped()
			return nil
		}
	}
}

// Stop stops the watcher.
func (cw *CertFileWatcher) Stop() error {
	if !cw.CanStop() {
		return async.ErrCannotStop
	}

	cw.Stopping()
	<-cw.NotifyStopped()
	return nil
}
