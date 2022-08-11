//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package server

func (s *Server) routes() {
	api := s.Group("/api")
	{
		user := api.Group("/user")
		{
			user.POST("/add", s.addUser())
			user.POST("/login", s.loginUser())
		}
		translate := api.Group("/translate").Use(s.auth())
		{
			translate.GET("/ping", s.pong())
		}
	}

}
