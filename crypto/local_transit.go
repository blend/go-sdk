/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package crypto

import (
	"crypto/aes"
	"io"
	"strings"

	"github.com/blend/go-sdk/ex"
)

// Important constants.
const (
	// DefaultKeySize is the size of keys to generate for client use.
	DefaultKeySize = 32
	// KeyVersionSize is the size of the key version prefix.
	KeyVersionSize = (4 + 2 + 2 + 1) // YYYY + MM + DD + :
	// IVSize is the size of the IV prefix.
	IVSize = aes.BlockSize
	// HashSize is the size of the hash prefix.
	HashSize = 32 // reasons.
)

// NewLocalTransit returns a new manager.
func NewLocalTransit(options ...LocalTransitOption) LocalTransit {
	var lt LocalTransit
	for _, opt := range options {
		opt(&lt)
	}
	return lt
}

// LocalTransitOption is an option for local transit types.
type LocalTransitOption func(*LocalTransit)

// OptLocalTransitKey is a local transit option..
func OptLocalTransitKey(key []byte) LocalTransitOption {
	return func(lt *LocalTransit) {
		lt.KeyProvider = func(_ string) ([]byte, error) {
			return key, nil
		}
	}
}

// OptLocalTransitKeyProvider is a local transit option..
func OptLocalTransitKeyProvider(provider func(string) ([]byte, error)) LocalTransitOption {
	return func(lt *LocalTransit) {
		lt.KeyProvider = provider
	}
}

// OptLocalTransitContextProvider is a local transit option..
func OptLocalTransitContextProvider(provider func() string) LocalTransitOption {
	return func(lt *LocalTransit) {
		lt.ContextProvider = provider
	}
}

// LocalTransit is a local transit encrypter and decrypter.
type LocalTransit struct {
	ContextProvider func() string
	KeyProvider     func(string) ([]byte, error)
}

// GetKey returns a key of a specific version.
func (m LocalTransit) GetKey(context string) ([]byte, error) {
	if m.KeyProvider == nil {
		return nil, ex.New("local transit; key provider is unset")
	}
	return m.KeyProvider(context)
}

// Encrypt reads a source stream and encrypts it to the destination stream.
func (m LocalTransit) Encrypt(dst io.Writer, src io.Reader) error {
	var currentContext string
	if m.ContextProvider != nil {
		currentContext = m.ContextProvider()
	}
	encKey, err := m.GetKey(currentContext)
	if err != nil {
		return ex.New(err)
	}
	macKey, err := m.GetKey(currentContext)
	if err != nil {
		return ex.New(err)
	}

	wr, err := NewStreamEncrypter(encKey, macKey, src)
	if err != nil {
		return ex.New(err)
	}

	versionPrefix := currentContext + ":"
	_, err = io.WriteString(dst, versionPrefix)
	if err != nil {
		return ex.New(err)
	}

	streamMeta := wr.Meta()
	_, err = dst.Write(streamMeta.IV)
	if err != nil {
		return ex.New(err)
	}

	_, err = dst.Write(streamMeta.Hash)
	if err != nil {
		return ex.New(err)
	}

	_, err = io.Copy(dst, wr)
	if err != nil {
		return ex.New(err)
	}

	return nil
}

// Decrypt reads a source stream and encrypts it to the destination stream.
func (m LocalTransit) Decrypt(dst io.Writer, src io.Reader) error {
	keyVersion, err := m.readKeyVersionPrefix(src)
	if err != nil {
		return ex.New(err)
	}

	encKey, err := m.GetKey(keyVersion)
	if err != nil {
		return ex.New(err)
	}
	macKey, err := m.GetKey(keyVersion)
	if err != nil {
		return ex.New(err)
	}

	iv, hash, err := m.readIVHashPrefix(src)
	if err != nil {
		return ex.New(err)
	}

	r, err := NewStreamDecrypter(encKey, macKey, StreamMeta{IV: iv, Hash: hash}, src)
	if err != nil {
		return ex.New(err)
	}

	// check the seal
	if err = r.Authenticate(); err != nil {
		return err
	}

	_, err = io.Copy(dst, r)
	if err != nil {
		return ex.New(err)
	}

	return nil
}

func (m LocalTransit) readKeyVersionPrefix(src io.Reader) (string, error) {
	keyVersion := make([]byte, KeyVersionSize)
	read, err := io.ReadFull(src, keyVersion)
	if err != nil {
		return "", ex.New(err)
	}
	if read != KeyVersionSize {
		return "", ex.New("local transit; invalid key version; wrong length")
	}
	if !strings.HasSuffix(string(keyVersion), ":") {
		return "", ex.New("local transit; invalid key version; wrong terminator")
	}
	return strings.TrimSuffix(string(keyVersion), ":"), nil
}

func (m LocalTransit) readIVHashPrefix(src io.Reader) (iv, hash []byte, err error) {
	iv = make([]byte, IVSize)
	hash = make([]byte, HashSize)

	if _, err = io.ReadFull(src, iv); err != nil {
		err = ex.New(err)
		return
	}

	if _, err = io.ReadFull(src, hash); err != nil {
		err = ex.New(err)
		return
	}

	return
}
