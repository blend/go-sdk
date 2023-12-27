/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package vault

import (
	"context"
	"encoding/base64"
)

// Assert MockTransitClient implements TransitClient
var (
	_ Client        = MockTransitClient{}
	_ TransitClient = MockTransitClient{}
)

// MockTransitClient skips interactions with the vault for encryption/decryption
type MockTransitClient struct {
	Client
}

// Encrypt just returns the input in the mock
func (m MockTransitClient) Encrypt(ctx context.Context, key string, context, data []byte) (string, error) {
	return base64.StdEncoding.EncodeToString(data), nil
}

// TransitHMAC just returns the input
func (m MockTransitClient) TransitHMAC(ctx context.Context, key string, input []byte) ([]byte, error) {
	return input, nil
}

// Decrypt just returns the input in the mock
func (m MockTransitClient) Decrypt(ctx context.Context, key string, context []byte, ciphertext string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(ciphertext)
}

// BatchEncrypt just returns the input
func (m MockTransitClient) BatchEncrypt(ctx context.Context, key string, batchInput BatchTransitInput) ([]string, error) {
	var res []string
	for i := range batchInput.BatchTransitInputItems {
		data := batchInput.BatchTransitInputItems[i].Plaintext
		res = append(res, base64.StdEncoding.EncodeToString(data))
	}
	return res, nil
}

// BatchDecrypt just returns the input
func (m MockTransitClient) BatchDecrypt(ctx context.Context, key string, batchInput BatchTransitInput) ([][]byte, error) {
	var res [][]byte
	for i := range batchInput.BatchTransitInputItems {
		data, err := base64.StdEncoding.DecodeString(batchInput.BatchTransitInputItems[i].Ciphertext)
		if err != nil {
			return nil, err
		}
		res = append(res, data)
	}
	return res, nil
}
