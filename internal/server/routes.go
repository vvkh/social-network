package server

import (
	"github.com/go-chi/chi/v5"

	"github.com/vvkh/social-network/internal/routes/friends"
	"github.com/vvkh/social-network/internal/routes/index"
	"github.com/vvkh/social-network/internal/routes/login"
	"github.com/vvkh/social-network/internal/routes/register"
	"github.com/vvkh/social-network/internal/routes/user"
	"github.com/vvkh/social-network/internal/routes/users"
)

func (s *server) setupRoutes() {
	s.handler.Get("/", index.Handle())
	s.handler.Get("/login/", login.Handle())
	s.handler.Get("/register/", register.Handle())
	s.handler.Get("/friends/", friends.Handle())
	s.handler.Route("/users/", func(r chi.Router) {
		r.Get("/", users.Handle())
		r.Get("/{userID}", user.Handle())
	})

}
