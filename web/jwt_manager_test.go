/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/crypto"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/uuid"
)

func TestNewJWTManager(t *testing.T) {
	assert := assert.New(t)

	key := crypto.MustCreateKey(32)
	m := NewJWTManager(key)
	assert.NotNil(m.KeyProvider)

	stored, err := m.KeyProvider(nil)
	assert.Nil(err)
	assert.Equal(key, stored)
}

func TestNewJWTManagerClaims(t *testing.T) {
	assert := assert.New(t)

	key := crypto.MustCreateKey(32)
	m := NewJWTManager(key)

	session := &Session{
		SessionID:	uuid.V4().String(),
		BaseURL:	uuid.V4().String(),
		UserID:		uuid.V4().String(),
		CreatedUTC:	time.Date(2018, 9, 8, 12, 00, 0, 0, time.UTC),
		ExpiresUTC:	time.Date(2018, 9, 9, 12, 00, 0, 0, time.UTC),
	}

	claims := m.Claims(session)
	assert.Equal(session.SessionID, claims.Id)
	assert.Equal(session.BaseURL, claims.Audience)
	assert.Equal("go-web", claims.Issuer)
	assert.Equal(session.UserID, claims.Subject)
	assert.Equal(session.CreatedUTC, time.Unix(claims.IssuedAt, 0).In(time.UTC))
	assert.Equal(session.ExpiresUTC, time.Unix(claims.ExpiresAt, 0).In(time.UTC))
}

func TestNewJWTManagerFromClaims(t *testing.T) {
	assert := assert.New(t)

	key := crypto.MustCreateKey(32)
	m := NewJWTManager(key)

	claims := &jwt.StandardClaims{
		Id:		uuid.V4().String(),
		Audience:	uuid.V4().String(),
		Issuer:		"go-web",
		Subject:	uuid.V4().String(),
		IssuedAt:	time.Date(2018, 9, 8, 12, 00, 0, 0, time.UTC).Unix(),
		ExpiresAt:	time.Date(2018, 9, 9, 12, 00, 0, 0, time.UTC).Unix(),
	}

	session := m.FromClaims(claims)
	assert.Equal(session.SessionID, claims.Id)
	assert.Equal(session.BaseURL, claims.Audience)
	assert.Equal(session.UserID, claims.Subject)
	assert.Equal(session.CreatedUTC, time.Unix(claims.IssuedAt, 0).In(time.UTC))
	assert.Equal(session.ExpiresUTC, time.Unix(claims.ExpiresAt, 0).In(time.UTC))
}

func TestNewJWTManagerKeyFunc(t *testing.T) {
	assert := assert.New(t)

	key := crypto.MustCreateKey(32)
	m := NewJWTManager(key)

	_, err := m.KeyFunc(&jwt.Token{
		Claims: jwt.MapClaims{},
	})

	assert.True(ex.Is(ErrJWTNonstandardClaims, err))

	claims := &jwt.StandardClaims{
		Id:		uuid.V4().String(),
		Audience:	uuid.V4().String(),
		Issuer:		"go-web",
		Subject:	uuid.V4().String(),
		IssuedAt:	time.Date(2018, 9, 8, 12, 00, 0, 0, time.UTC).Unix(),
		ExpiresAt:	time.Date(2018, 9, 9, 12, 00, 0, 0, time.UTC).Unix(),
	}
	returnedKey, err := m.KeyFunc(&jwt.Token{
		Claims: claims,
	})
	assert.Nil(err)
	assert.Equal(key, returnedKey)
}

func TestNewJWTManagerSerialization(t *testing.T) {
	assert := assert.New(t)

	key := crypto.MustCreateKey(32)
	m := NewJWTManager(key)

	session := &Session{
		SessionID:	uuid.V4().String(),
		BaseURL:	uuid.V4().String(),
		UserID:		uuid.V4().String(),
		CreatedUTC:	time.Now().UTC(),
		ExpiresUTC:	time.Now().UTC().Add(time.Hour),
	}

	output, err := m.SerializeHandler(context.TODO(), session)
	assert.Nil(err)
	assert.NotEmpty(output)

	parsed, err := m.FetchHandler(context.TODO(), output)
	assert.Nil(err)
	assert.Equal(parsed.SessionID, session.SessionID)
	assert.Equal(parsed.BaseURL, session.BaseURL)
	assert.Equal(parsed.UserID, session.UserID)
	assert.False(parsed.CreatedUTC.IsZero())
	assert.False(parsed.ExpiresUTC.IsZero())
}
