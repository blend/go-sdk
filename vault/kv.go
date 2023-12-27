/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package vault

import "context"

// KV is a basic key value store.
type KV interface {
	Put(ctx context.Context, path string, data Values, options ...CallOption) error
	Get(ctx context.Context, path string, options ...CallOption) (Values, error)
	Delete(ctx context.Context, path string, options ...CallOption) error
	List(ctx context.Context, path string, options ...CallOption) ([]string, error)
}

// KVClient is a basic key value store client.
type KVClient = KV
