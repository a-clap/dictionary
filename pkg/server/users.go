//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package server

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"time"
)

var jwtTestKey []byte = []byte("test key")

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	claims   struct {
		Name string `json:"name"`
		jwt.RegisteredClaims
	}
}

var (
	ErrExpired = errors.New("token expired")
)

func (s *Server) addUser() gin.HandlerFunc {
	return func(context *gin.Context) {
		var user User
		if err := context.ShouldBindJSON(&user); err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := s.u.Add(user.Name, user.Password); err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		context.JSON(http.StatusCreated, gin.H{"name": user.Name})
	}
}

func (s *Server) loginUser() gin.HandlerFunc {
	return func(context *gin.Context) {
		var user User
		if err := context.ShouldBindJSON(&user); err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if login, err := s.u.Auth(user.Name, user.Password); err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else if !login {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		// All good, generate token
		token, err := user.Token(30 * time.Minute)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusOK, gin.H{"token": token})
	}
}

func (u *User) Token(duration time.Duration) (string, error) {
	expires := time.Now().Add(duration)
	u.claims.Name = u.Name
	u.claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{Time: expires},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, u.claims)
	return token.SignedString(jwtTestKey)
}

func (u *User) Validate(token string) (bool, error) {
	tkn, err := jwt.ParseWithClaims(token, &u.claims, func(token *jwt.Token) (interface{}, error) {
		return jwtTestKey, nil
	})

	if err != nil {
		if validationError, ok := err.(*jwt.ValidationError); ok {
			if (validationError.Errors & jwt.ValidationErrorExpired) == jwt.ValidationErrorExpired {
				return false, ErrExpired
			}
		}
		return false, err
	}

	return tkn.Valid, nil
}
