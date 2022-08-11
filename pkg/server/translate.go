//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) pong() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"message": "pong"})
	}
}
