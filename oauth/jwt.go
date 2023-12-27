/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package oauth

import (
	"golang.org/x/oauth2"

	"github.com/golang-jwt/jwt/v4"

	"github.com/blend/go-sdk/ex"
)

// ParseTokenJWT parses a jwt from a given oauth2 token.
func ParseTokenJWT(tok *oauth2.Token, keyfunc jwt.Keyfunc) (*GoogleClaims, error) {
	jwtRaw, ok := tok.Extra("id_token").(string)
	if !ok || jwtRaw == "" {
		return nil, ex.New("invalid oauth token; `id_token` jwt missing")
	}
	var claims GoogleClaims
	_, err := jwt.ParseWithClaims(jwtRaw, &claims, keyfunc)
	if err != nil {
		return nil, err
	}
	return &claims, nil
}

// GoogleClaims are extensions to the jwt standard claims for google oauth.
//
// See additional documentation here: https://developers.google.com/identity/sign-in/web/backend-auth
type GoogleClaims struct {
	jwt.StandardClaims

	Email		string	`json:"email"`
	EmailVerified	string	`json:"email-verified"`
	HD		string	`json:"hd"`
	Nonce		string	`json:"nonce"`

	FamilyName	string	`json:"family_name"`
	GivenName	string	`json:"given_name"`
	Locale		string	`json:"locale"`
	Picture		string	`json:"picture"`
	Profile		string	`json:"profile"`
}
