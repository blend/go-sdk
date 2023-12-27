/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package proxyprotocol

import (
	"crypto/tls"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestCreateListenerOptions(t *testing.T) {
	assert := assert.New(t)

	var options CreateListenerOptions

	assert.False(options.KeepAlive)
	assert.Nil(OptKeepAlive(true)(&options))
	assert.True(options.KeepAlive)

	assert.Zero(options.KeepAlivePeriod)
	assert.Nil(OptKeepAlivePeriod(time.Second)(&options))
	assert.Equal(time.Second, options.KeepAlivePeriod)

	assert.False(options.UseProxyProtocol)
	assert.Nil(OptUseProxyProtocol(true)(&options))
	assert.True(options.UseProxyProtocol)

	assert.Nil(options.TLSConfig)
	assert.Nil(OptTLSConfig(&tls.Config{})(&options))
	assert.NotNil(options.TLSConfig)
}
