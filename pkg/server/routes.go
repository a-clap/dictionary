//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package server

func (s *Server) routes() {
	s.POST("/api/user/add", s.addUser())
	s.POST("/api/user/login", s.loginUser())
}
