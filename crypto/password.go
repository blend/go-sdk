/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package crypto

import "golang.org/x/crypto/bcrypt"

const (
	bcryptHashingCost = 10
)

// HashPassword uses bcrypt to generate a salted hash for the provided password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptHashingCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// PasswordMatchesHash checks whether the provided password matches the provided hash
func PasswordMatchesHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
