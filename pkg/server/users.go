//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package server

import "github.com/gin-gonic/gin"

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (s *Server) addUser() gin.HandlerFunc {
	return func(context *gin.Context) {

	}
}
