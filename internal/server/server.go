package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	*gin.Engine
}

func (s *Server) Translate() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "%s", "mózg")
	}
}

func New() *Server {
	s := &Server{gin.Default()}
	s.routes()
	return s
}
