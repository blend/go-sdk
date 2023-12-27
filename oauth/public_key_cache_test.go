/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package oauth

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/jwk"
	"github.com/blend/go-sdk/r2"
	"github.com/blend/go-sdk/uuid"
)

func Test_PublicKeyCache_Keyfunc(t *testing.T) {
	it := assert.New(t)

	pk0, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(pk0pem))
	it.Nil(err)
	pk1, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(pk1pem))
	it.Nil(err)
	keys := []jwk.JWK{
		createJWK(pk0),
		createJWK(pk1),
	}

	var didCallResponder bool
	keysResponder := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		didCallResponder = true
	}))
	defer keysResponder.Close()

	cache := NewPublicKeyCache(GoogleKeysURL)
	cache.FetchPublicKeysDefaults = []r2.Option{
		r2.OptURL(keysResponder.URL),
	}
	cache.current = &PublicKeysResponse{
		CacheControl:	"public, max-age=23196, must-revalidate, no-transform",
		Expires:	time.Now().UTC().AddDate(0, 0, 1),
		Keys:		jwkLookup(keys),
	}

	keyfunc := cache.Keyfunc(context.TODO())

	pub, err := keyfunc(&jwt.Token{
		Header: map[string]interface{}{
			"kid": keys[0].KID,
		},
	})

	it.Nil(err)

	typedPub, ok := pub.(*rsa.PublicKey)
	it.True(ok)
	it.Equal(*pk0.PublicKey.N, *typedPub.N)
	it.False(didCallResponder)
}

func Test_PublicKeyCache_Keyfunc_MissingKIDHeader(t *testing.T) {
	it := assert.New(t)

	cache := NewPublicKeyCache(GoogleKeysURL)
	keyfunc := cache.Keyfunc(context.TODO())
	pub, err := keyfunc(&jwt.Token{})
	it.NotNil(err)
	it.Nil(pub)
}

func Test_PublicKeyCache_Keyfunc_InvalidKID(t *testing.T) {
	it := assert.New(t)

	cache := NewPublicKeyCache(GoogleKeysURL)
	keyfunc := cache.Keyfunc(context.TODO())
	pub, err := keyfunc(&jwt.Token{
		Header: map[string]interface{}{
			"kid": 1234,
		},
	})
	it.NotNil(err)
	it.Nil(pub)
}

func Test_PublicKeyCache_Get(t *testing.T) {
	it := assert.New(t)

	pk0, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(pk0pem))
	it.Nil(err)
	pk1, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(pk1pem))
	it.Nil(err)
	keys := []jwk.JWK{
		createJWK(pk0),
		createJWK(pk1),
	}
	var didCallResponder bool
	keysResponder := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		didCallResponder = true
		rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
		rw.Header().Set("Cache-Control", "public, max-age=23196, must-revalidate, no-transform")	// set cache control
		rw.Header().Set("Expires", time.Now().UTC().AddDate(0, 1, 0).Format(http.TimeFormat))		// set expires
		rw.Header().Set("Date", time.Now().UTC().Format(http.TimeFormat))				// set date
		rw.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(rw).Encode(struct {
			Keys []jwk.JWK `json:"keys"`
		}{
			Keys: keys,
		})
	}))
	defer keysResponder.Close()

	cache := NewPublicKeyCache(GoogleKeysURL)
	cache.FetchPublicKeysDefaults = []r2.Option{
		r2.OptURL(keysResponder.URL),
	}

	pub, err := cache.Get(context.TODO(), keys[0].KID)
	it.Nil(err)
	it.NotNil(pub)
	it.Equal(*pk0.PublicKey.N, *pub.N)
	it.True(didCallResponder)

	didCallResponder = false

	pub, err = cache.Get(context.TODO(), keys[1].KID)
	it.Nil(err)
	it.NotNil(pub)
	it.Equal(*pk1.PublicKey.N, *pub.N)
	it.False(didCallResponder, "we should have cached the results")
}

func Test_PublicKeyCache_Get_NoRefresh(t *testing.T) {
	it := assert.New(t)

	pk0, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(pk0pem))
	it.Nil(err)
	pk1, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(pk1pem))
	it.Nil(err)
	keys := []jwk.JWK{
		createJWK(pk0),
		createJWK(pk1),
	}
	var didCallResponder bool
	keysResponder := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		didCallResponder = true
	}))
	defer keysResponder.Close()

	cache := NewPublicKeyCache(GoogleKeysURL)
	cache.current = &PublicKeysResponse{
		CacheControl:	"public, max-age=23196, must-revalidate, no-transform",
		Expires:	time.Now().UTC().AddDate(0, 0, 1),
		Keys:		jwkLookup(keys),
	}
	cache.FetchPublicKeysDefaults = []r2.Option{
		r2.OptURL(keysResponder.URL),
	}

	pub, err := cache.Get(context.TODO(), keys[0].KID)
	it.Nil(err)
	it.NotNil(pub)
	it.Equal(*pk0.PublicKey.N, *pub.N)
	it.False(didCallResponder)
}

func Test_PublicKeyCache_Get_Refresh(t *testing.T) {
	it := assert.New(t)

	pk0, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(pk0pem))
	it.Nil(err)
	pk1, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(pk1pem))
	it.Nil(err)
	keys := []jwk.JWK{
		createJWK(pk0),
		createJWK(pk1),
	}
	var didCallResponder bool
	keysResponder := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		didCallResponder = true
		rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
		rw.Header().Set("Cache-Control", "public, max-age=23196, must-revalidate, no-transform")	// set cache control
		rw.Header().Set("Expires", time.Now().UTC().AddDate(0, 1, 0).Format(http.TimeFormat))		// set expires
		rw.Header().Set("Date", time.Now().UTC().Format(http.TimeFormat))				// set date
		rw.WriteHeader(200)
		_ = json.NewEncoder(rw).Encode(struct {
			Keys []jwk.JWK `json:"keys"`
		}{
			Keys: keys,
		})
	}))
	defer keysResponder.Close()

	cache := NewPublicKeyCache(GoogleKeysURL)
	cache.current = &PublicKeysResponse{
		CacheControl:	"public, max-age=23196, must-revalidate, no-transform",
		Expires:	time.Now().UTC().AddDate(0, 0, -1),
		Keys:		jwkLookup(keys),
	}
	cache.FetchPublicKeysDefaults = []r2.Option{
		r2.OptURL(keysResponder.URL),
	}

	pub, err := cache.Get(context.TODO(), keys[0].KID)
	it.Nil(err)
	it.NotNil(pub)
	it.Equal(*pk0.PublicKey.N, *pub.N)
	it.True(didCallResponder)
}

func Test_PublicKeyCache_Get_RefreshOnMiss(t *testing.T) {
	it := assert.New(t)

	pk0, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(pk0pem))
	it.Nil(err)
	pk1, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(pk1pem))
	it.Nil(err)
	keys := []jwk.JWK{
		createJWK(pk0),
		createJWK(pk1),
	}
	var didCallResponder bool
	keysResponder := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		didCallResponder = true
		rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
		rw.Header().Set("Cache-Control", "public, max-age=23196, must-revalidate, no-transform")	// set cache control
		rw.Header().Set("Expires", time.Now().UTC().AddDate(0, 1, 0).Format(http.TimeFormat))		// set expires
		rw.Header().Set("Date", time.Now().UTC().Format(http.TimeFormat))				// set date
		rw.WriteHeader(200)
		_ = json.NewEncoder(rw).Encode(struct {
			Keys []jwk.JWK `json:"keys"`
		}{
			Keys: keys,
		})
	}))
	defer keysResponder.Close()

	cache := NewPublicKeyCache(GoogleKeysURL)
	cache.current = &PublicKeysResponse{
		CacheControl:	"public, max-age=23196, must-revalidate, no-transform",
		Expires:	time.Now().UTC().AddDate(0, 0, -1),
		Keys:		jwkLookup(keys),
	}
	cache.FetchPublicKeysDefaults = []r2.Option{
		r2.OptURL(keysResponder.URL),
	}

	pub, err := cache.Get(context.TODO(), uuid.V4().String())
	it.NotNil(err)
	it.Nil(pub)
	it.True(didCallResponder)
}
