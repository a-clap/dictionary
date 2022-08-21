//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package server

import (
	"github.com/a-clap/dictionary/internal/auth"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	auth.StoreTokener
}

type Server struct {
	*gin.Engine
	manager *auth.Manager
}

func New(h Handler) *Server {
	s := &Server{
		Engine:  gin.Default(),
		manager: auth.New(h),
	}

	s.routes()
	return s
}
