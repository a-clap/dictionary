//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package server

import (
	"github.com/a-clap/dictionary/internal/users"
	"github.com/gin-gonic/gin"
)

type Server struct {
	*gin.Engine
	u *users.Users
}

func New(u *users.Users) *Server {
	return &Server{
		Engine: gin.Default(),
		u:      u,
	}
}
