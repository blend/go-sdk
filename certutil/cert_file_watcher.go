/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

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

const (
	// DefaultCertficicateFileWatcherPollInterval is the default poll interval when re-reading certs
	DefaultCertficicateFileWatcherPollInterval = 500 * time.Millisecond
)

// NewCertFileWatcher creates a new CertReloader object with a reload delay
func NewCertFileWatcher(keyPair KeyPair, opts ...CertFileWatcherOption) (*CertFileWatcher, error) {
	if keyPair.CertPath == "" || keyPair.KeyPath == "" {
		return nil, ex.New(ErrTLSPathsUnset)
	}
	cw := &CertFileWatcher{
		latch:		async.NewLatch(),
		keyPair:	keyPair,
	}
	for _, opt := range opts {
		if err := opt(cw); err != nil {
			return nil, err
		}
	}
	cert, err := tls.LoadX509KeyPair(cw.keyPair.CertPath, cw.keyPair.KeyPath)
	if err != nil {
		return nil, err
	}
	cw.certificate = &cert
	return cw, nil
}

// CertFileWatcherOption is an option for a cert watcher.
type CertFileWatcherOption func(*CertFileWatcher) error

// CertFileWatcherOnReloadAction is the on reload action for a cert file watcher.
type CertFileWatcherOnReloadAction func(*CertFileWatcher) error

// OptCertFileWatcherOnReload sets the on reload handler.
// If you need to capture *every* reload of the cert, including the initial one in the constructor
// you must use this option.
func OptCertFileWatcherOnReload(handler CertFileWatcherOnReloadAction) CertFileWatcherOption {
	return func(cfw *CertFileWatcher) error {
		cfw.onReload = handler
		return nil
	}
}

// OptCertFileWatcherNotifyReload sets the notify reload channel.
func OptCertFileWatcherNotifyReload(notifyReload chan struct{}) CertFileWatcherOption {
	return func(cfw *CertFileWatcher) error {
		cfw.notifyReload = notifyReload
		return nil
	}
}

// OptCertFileWatcherPollInterval sets the poll interval .
func OptCertFileWatcherPollInterval(d time.Duration) CertFileWatcherOption {
	return func(cfw *CertFileWatcher) error {
		cfw.pollInterval = d
		return nil
	}
}

// CertFileWatcher reloads a cert key pair when there is a change, e.g. cert renewal
type CertFileWatcher struct {
	latch		*async.Latch
	certificateMu	sync.RWMutex
	certificate	*tls.Certificate
	keyPair		KeyPair
	pollInterval	time.Duration
	notifyReload	chan struct{}
	onReload	CertFileWatcherOnReloadAction
}

// CertPath returns the cert path.
func (cw *CertFileWatcher) CertPath() string	{ return cw.keyPair.CertPath }

// KeyPath returns the cert path.
func (cw *CertFileWatcher) KeyPath() string	{ return cw.keyPair.KeyPath }

// PollIntervalOrDefault returns the polling interval or a default.
func (cw *CertFileWatcher) PollIntervalOrDefault() time.Duration {
	if cw.pollInterval > 0 {
		return cw.pollInterval
	}
	return DefaultCertficicateFileWatcherPollInterval
}

// Reload forces the reload of the underlying certificate.
func (cw *CertFileWatcher) Reload() (err error) {
	defer func() {
		if cw.notifyReload != nil {
			cw.notifyReload <- struct{}{}
		}
		if cw.onReload != nil && err == nil {
			err = cw.onReload(cw)
		}
	}()

	cert, loadErr := tls.LoadX509KeyPair(cw.keyPair.CertPath, cw.keyPair.KeyPath)
	if loadErr != nil {
		err = ex.New(loadErr)
		return
	}
	cw.certificateMu.Lock()
	cw.certificate = &cert
	cw.certificateMu.Unlock()
	return
}

// Certificate gets the underlying certificate, it blocks when the `cert` field is being updated
func (cw *CertFileWatcher) Certificate() *tls.Certificate {
	cw.certificateMu.RLock()
	defer cw.certificateMu.RUnlock()
	return cw.certificate
}

// GetCertificate gets the underlying certificate in the form that tls config expects.
func (cw *CertFileWatcher) GetCertificate(_ *tls.ClientHelloInfo) (*tls.Certificate, error) {
	cw.certificateMu.RLock()
	defer cw.certificateMu.RUnlock()
	return cw.certificate, nil
}

// IsStarted returns if the underlying latch is started.
func (cw *CertFileWatcher) IsStarted() bool	{ return cw.latch.IsStarted() }

// IsStopped returns if the underlying latch is stopped.
func (cw *CertFileWatcher) IsStopped() bool	{ return cw.latch.IsStopped() }

// NotifyStarted returns the notify started channel.
func (cw *CertFileWatcher) NotifyStarted() <-chan struct{} {
	return cw.latch.NotifyStarted()
}

// NotifyStopped returns the notify stopped channel.
func (cw *CertFileWatcher) NotifyStopped() <-chan struct{} {
	return cw.latch.NotifyStopped()
}

// NotifyReload the notify reload channel.
//
// You must supply this channel as an option in the constructor.
func (cw *CertFileWatcher) NotifyReload() <-chan struct{} {
	return cw.notifyReload
}

// Start watches the cert and triggers a reload on change
func (cw *CertFileWatcher) Start() error {
	cw.latch.Starting()

	certLastMod, keyLastMod, err := cw.keyPairLastModified()
	if err != nil {
		cw.latch.Stopped()
		return err
	}

	ticker := time.NewTicker(cw.PollIntervalOrDefault())
	defer ticker.Stop()

	cw.latch.Started()
	var certMod, keyMod time.Time
	for {
		select {
		case <-ticker.C:
			certMod, keyMod, err = cw.keyPairLastModified()
			if err != nil {
				return err
			}
			// wait for both to update
			if keyMod.After(keyLastMod) && certMod.After(certLastMod) {
				if err = cw.Reload(); err != nil {
					return err
				}
				keyLastMod = keyMod
				certLastMod = certMod
			}
		case <-cw.latch.NotifyStopping():
			cw.latch.Stopped()
			return nil
		}
	}
}

// Stop stops the watcher.
func (cw *CertFileWatcher) Stop() error {
	if !cw.latch.CanStop() {
		return async.ErrCannotStop
	}
	cw.latch.WaitStopped()
	cw.latch.Reset()
	return nil
}

func (cw *CertFileWatcher) keyPairLastModified() (cert time.Time, key time.Time, err error) {
	var certStat, keyStat os.FileInfo
	certStat, err = os.Stat(cw.keyPair.CertPath)
	if err != nil {
		return
	}
	keyStat, err = os.Stat(cw.keyPair.KeyPath)
	if err != nil {
		return
	}
	cert = certStat.ModTime()
	key = keyStat.ModTime()
	return
}
