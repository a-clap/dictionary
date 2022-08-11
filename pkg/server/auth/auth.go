//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Claims struct {
	Name string `json:"name"`
	jwt.RegisteredClaims
}

var jwtTestKey []byte = []byte("some key")

func Token(name string, duration time.Duration) (string, error) {
	expires := time.Now().Add(duration)
	claims := &Claims{
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: expires},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtTestKey)
}

func Validate(token string) (bool, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtTestKey, nil
	})

	if err != nil {
		return false, err
	}

	if claims.ExpiresAt.Unix() < time.Now().Local().Unix() {
		return false, nil
	}

	return tkn.Valid, nil
}
