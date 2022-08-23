//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package server

import (
	"github.com/a-clap/dictionary/internal/auth"
	"github.com/a-clap/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		logger.Infof("auth")
		token := context.GetHeader("Authorization")
		if len(token) == 0 {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "request doesn't contain an authorization token"})
			return
		}
		user, err := s.manager.ValidateToken(token)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		logger.Infof("user %s logged successfully", user.Name)
		context.Next()
	}
}

func (s *Server) addUser() gin.HandlerFunc {
	return func(context *gin.Context) {
		var user auth.User
		if err := context.ShouldBindJSON(&user); err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := s.manager.Add(user); err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		context.JSON(http.StatusCreated, gin.H{"name": user.Name})
	}
}

func (s *Server) loginUser() gin.HandlerFunc {
	return func(context *gin.Context) {
		var user auth.User
		if err := context.ShouldBindJSON(&user); err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if login, err := s.manager.Auth(user); err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else if !login {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		// All good, generate token
		token, err := s.manager.Token(user)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusOK, gin.H{"token": token})
	}
}
