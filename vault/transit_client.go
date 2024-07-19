/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package vault

import (
	"context"
)

// TransitClient is an interface for an encryption-as-a-service client
type TransitClient interface {
	CreateTransitKey(ctx context.Context, key string, options ...CreateTransitKeyOption) error
	ConfigureTransitKey(ctx context.Context, key string, options ...UpdateTransitKeyOption) error
	ReadTransitKey(ctx context.Context, key string) (map[string]interface{}, error)
	DeleteTransitKey(ctx context.Context, key string) error

	Encrypt(ctx context.Context, key string, context, data []byte) (string, error)
	Decrypt(ctx context.Context, key string, context []byte, ciphertext string) ([]byte, error)

	TransitHMAC(ctx context.Context, key string, input []byte) ([]byte, error)

	BatchEncrypt(ctx context.Context, key string, batchInput BatchTransitInput) ([]string, error)
	BatchDecrypt(ctx context.Context, key string, batchInput BatchTransitInput) ([][]byte, error)
}
