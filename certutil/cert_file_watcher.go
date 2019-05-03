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

	OnReload func(*tls.Certificate, error)
}

// PollIntervalOrDefault returns the polling interval or a default.
func (cw *CertFileWatcher) PollIntervalOrDefault() time.Duration {
	if cw.PollInterval > 0 {
		return cw.PollInterval
	}
	return 500 * time.Millisecond
}

// Reload forces the reload of the underlying certificate.
func (cw *CertFileWatcher) Reload() (err error) {
	cw.syncRoot.Lock()
	defer cw.syncRoot.Unlock()
	defer func() {
		if cw.OnReload != nil {
			cw.OnReload(cw.Certificate, err)
		}
	}()

	cert, loadErr := tls.LoadX509KeyPair(cw.CertPath, cw.KeyPath)
	if loadErr != nil {
		err = ex.New(loadErr)
		return
	}
	cw.Certificate = &cert
	return
}

// GetCertificate gets the cached certificate, it blocks when the `cert` field is being updated
func (cw *CertFileWatcher) GetCertificate(_ *tls.ClientHelloInfo) (*tls.Certificate, error) {
	cw.syncRoot.RLock()
	defer cw.syncRoot.RUnlock()
	return cw.Certificate, nil
}

// Start watches the cert and triggers a reload on change
func (cw *CertFileWatcher) Start() error {
	cw.Starting()

	certLastMod, keyLastMod, err := cw.keyPairLastModified()
	if err != nil {
		return err
	}

	ticker := time.Tick(cw.PollIntervalOrDefault())
	cw.Started()
	var certMod, keyMod time.Time
	for {
		select {
		case <-ticker:
			certMod, keyMod, err = cw.keyPairLastModified()
			if err != nil {
				return err
			}
			if keyMod.After(keyLastMod) || certMod.After(certLastMod) {
				if err = cw.Reload(); err != nil {
					return err
				}
				keyLastMod = keyMod
				certLastMod = certMod
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

func (cw *CertFileWatcher) keyPairLastModified() (cert time.Time, key time.Time, err error) {
	var certStat, keyStat os.FileInfo
	keyStat, err = os.Stat(cw.KeyPath)
	if err != nil {
		return
	}
	certStat, err = os.Stat(cw.CertPath)
	if err != nil {
		return
	}

	cert = certStat.ModTime()
	key = keyStat.ModTime()
	return
}
