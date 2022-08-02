package server

func (s *Server) routes() {
	s.GET("/translate/en/brain", s.Translate())
}
