//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package server

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
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
		token, err := user.token()
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusOK, gin.H{"token": token})
	}
}

var jwtTestKey []byte = []byte("some key")

func (u *User) token() (string, error) {
	expires := time.Now().Add(30 * time.Minute)
	claims := &Claims{
		Name: u.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: expires},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SigningString()
}

func (u *User) validateToken(token string) (bool, error) {
	return false, nil
}
