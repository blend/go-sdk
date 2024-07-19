/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package vault

import "github.com/blend/go-sdk/ex"

var validTKTypes = map[string]struct{}{
	TypeAES256GCM96:      {},
	TypeCHACHA20POLY1305: {},
	TypeED25519:          {},
	TypeECDSAP256:        {},
	TypeRSA2048:          {},
	TypeRSA4096:          {},
}

// CreateTransitKeyOption is an option type for transit key creation
type CreateTransitKeyOption func(tkc *CreateTransitKeyConfig) error

// UpdateTransitKeyOption is an option type for transit key creation
type UpdateTransitKeyOption func(tkc *UpdateTransitKeyConfig) error

// OptCreateTransitConfig is a creation option for when you have a pre-defined struct
func OptCreateTransitConfig(config CreateTransitKeyConfig) CreateTransitKeyOption {
	return func(tkc *CreateTransitKeyConfig) error {
		*tkc = config
		return nil
	}
}

// OptCreateTransitConvergent - If enabled, the key will support convergent encryption, where the same plaintext creates
// the same ciphertext. This also sets derived to true (which is required). When enabled, each encryption (or decryption
// or rewrap or datakey) operation will derive a nonce value rather than randomly generate it.
func OptCreateTransitConvergent() CreateTransitKeyOption {
	return func(tkc *CreateTransitKeyConfig) error {
		tkc.Convergent = true
		tkc.Derived = true
		return nil
	}
}

// OptCreateTransitDerived - Specifies if key derivation is to be used. If enabled, all encrypt/decrypt requests to this
// named key must provide a context which is used for key derivation.
func OptCreateTransitDerived() CreateTransitKeyOption {
	return func(tkc *CreateTransitKeyConfig) error {
		tkc.Derived = true
		return nil
	}
}

// OptCreateTransitExportable - Enables keys to be exportable. This allows for all the valid keys in the key ring to be
// exported. Once set, this cannot be disabled.
func OptCreateTransitExportable() CreateTransitKeyOption {
	return func(tkc *CreateTransitKeyConfig) error {
		tkc.Exportable = true
		return nil
	}
}

// OptCreateTransitAllowPlaintextBackup - If set, enables taking backup of named key in the plaintext format. Once set, this
// cannot be disabled.
func OptCreateTransitAllowPlaintextBackup() CreateTransitKeyOption {
	return func(tkc *CreateTransitKeyConfig) error {
		tkc.AllowPlaintextBackup = true
		return nil
	}
}

// OptCreateTransitType - specifies the type of key to create. The default type is "aes256-gcm96":
//
//	  aes256-gcm96 – AES-256 wrapped with GCM using a 96-bit nonce size AEAD (symmetric, supports derivation and
//	     convergent encryption)
//	  chacha20-poly1305 – ChaCha20-Poly1305 AEAD (symmetric, supports derivation and convergent encryption)
//	  ed25519 – ED25519 (asymmetric, supports derivation). When using derivation, a sign operation with the same
//	     context will derive the same key and signature; this is a signing analog to convergent_encryption.
//		 ecdsa-p256 – ECDSA using the P-256 elliptic curve (asymmetric)
//		 rsa-2048 - RSA with bit size of 2048 (asymmetric)
//	  rsa-4096 - RSA with bit size of 4096 (asymmetric)
func OptCreateTransitType(keyType string) CreateTransitKeyOption {
	return func(tkc *CreateTransitKeyConfig) error {
		if _, ok := validTKTypes[keyType]; !ok {
			return ex.New("invalid keyType")
		}
		tkc.Type = keyType
		return nil
	}
}

// OptUpdateTransitConfig is an update option for when you have a pre-defined struct
func OptUpdateTransitConfig(config UpdateTransitKeyConfig) UpdateTransitKeyOption {
	return func(tku *UpdateTransitKeyConfig) error {
		*tku = config
		return nil
	}
}

// OptUpdateTransitMinDecryptionVer - Specifies the minimum version of ciphertext allowed to be decrypted. Adjusting
// this as part of a key rotation policy can prevent old copies of ciphertext from being decrypted, should they fall into
// the wrong hands. For signatures, this value controls the minimum version of signature that can be verified
// against. For HMACs, this controls the minimum version of a key allowed to be used as the key for verification.
func OptUpdateTransitMinDecryptionVer(minDecryptionVersion int) UpdateTransitKeyOption {
	return func(tku *UpdateTransitKeyConfig) error {
		tku.MinDecryptionVersion = minDecryptionVersion
		return nil
	}
}

// OptUpdateTransitMinEncryptionVer - Specifies the minimum version of the key that can be used to encrypt plaintext,
// sign payloads, or generate HMACs. Must be 0 (which will use the latest version) or a value greater or equal to
// min_decryption_version.
func OptUpdateTransitMinEncryptionVer(minEncryptionVersion int) UpdateTransitKeyOption {
	return func(tku *UpdateTransitKeyConfig) error {
		tku.MinEncryptionVersion = minEncryptionVersion
		return nil
	}
}

// OptUpdateTransitDeletionAllowed - Specifies if the key is allowed to be deleted.
func OptUpdateTransitDeletionAllowed(deletionAllowed bool) UpdateTransitKeyOption {
	return func(tku *UpdateTransitKeyConfig) error {
		tku.DeletionAllowed = &deletionAllowed
		return nil
	}
}

// OptUpdateTransitExportable - Enables keys to be exportable. This allows for all the valid keys in the key ring to be
// exported. Once set, this cannot be disabled.
func OptUpdateTransitExportable() UpdateTransitKeyOption {
	return func(tku *UpdateTransitKeyConfig) error {
		tku.Exportable = true
		return nil
	}
}

// OptUpdateTransitAllowPlaintextBackup - If set, enables taking backup of named key in the plaintext format. Once set, this
// cannot be disabled.
func OptUpdateTransitAllowPlaintextBackup() UpdateTransitKeyOption {
	return func(tku *UpdateTransitKeyConfig) error {
		tku.AllowPlaintextBackup = true
		return nil
	}
}
