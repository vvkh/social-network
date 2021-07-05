package server

import "github.com/vvkh/social-network/internal/routes/index"

func (s *server) setupRoutes() {
	s.handler.HandleFunc("/", index.HandleIndex())
}
