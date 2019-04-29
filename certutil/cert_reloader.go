package certutil

import (
	"crypto/tls"

	"github.com/blend/go-sdk/async"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
)

// NewCertReloader creates a new CertReloader object with a reload delay
func NewCertReloader(certPath, keyPath string) (*CertReloader, error) {
	result := &CertReloader{

		CertPath: certPath,
		KeyPath:  keyPath,
	}

	// load cert to make sure the current key pair is valid
	if err := result.Reload(); err != nil {
		return nil, err
	}

	return result, nil
}

// CertReloader reloads a cert key pair when there is a change, e.g. cert renewal
type CertReloader struct {
	*async.Latch

	Log         logger.Log
	Certificate *tls.Certificate

	CertPath string
	KeyPath  string
}

// Reload forces the reload of the underlying certificate.
func (cr *CertReloader) Reload() error {
	cr.Lock()
	defer cr.Unlock()

	cert, err := tls.LoadX509KeyPair(cr.CertPath, cr.KeyPath)
	if err != nil {
		return ex.New(err)
	}
	cr.Certificate = &cert
	return nil
}

// GetCertificate gets the cached certificate, it blocks when the `cert` field is being updated
func (cr *CertReloader) GetCertificate(_ *tls.ClientHelloInfo) (*tls.Certificate, error) {
	cr.mutex.RLock()
	defer cr.mutex.RUnlock()
	return cr.cert, nil
}

// State returns the current cert reloader state
func (cr *CertReloader) State() CertReloaderState {
	return cr.state
}

// Run watches the cert and triggers a reload on change
func (cr *CertReloader) Run() error {
	// we want to add the watcher and start polling for event right away (vs adding it in the constructor)
	if err := cr.watcher.Add(cr.certPath); err != nil {
		return exception.New(err)
	}
	defer cr.watcher.Remove(cr.certPath)
	cr.log.Infof("watching cert at %s", cr.certPath)

	cr.state = CertReloaderStateRunning
	defer func() {
		cr.state = CertReloaderStateStopped
	}()

	for {
		select {
		case event, ok := <-cr.watcher.EventsChan():
			if !ok {
				return nil
			}
			cr.log.Debugf("fsnotify event: %v", event)

			// note: behavior observed when kube updates mounted secret is chmod, remove, then the watch is lost
			// the behavior seems to vary across environment (e.g. kube, jenkins, mac)
			// if this causes any more issue, we should consider just doing poll + stat
			modified := fsnotify.Write | fsnotify.Remove | fsnotify.Chmod
			if event.Op&modified > 0 {
				// since we are watching only the cert, we may need to wait a bit for the key to get updated.
				cr.log.Infof("cert modified: %s, scheduling a reload after %s", event.Name, cr.reloadDelay)
				cr.scheduleReload()
			}
		case err, ok := <-cr.watcher.ErrorsChan():
			if !ok {
				return nil
			}
			cr.log.Errorf("fsnotify error: %v", err)
		case <-cr.stop:
			cr.log.Infof("stopped watching cert at %s", cr.certPath)
			return nil
		}
	}
}
