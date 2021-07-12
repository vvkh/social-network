package server

import (
	"github.com/go-chi/chi/v5"

	"github.com/vvkh/social-network/internal/routes/friends"
	"github.com/vvkh/social-network/internal/routes/index"
	"github.com/vvkh/social-network/internal/routes/login"
	"github.com/vvkh/social-network/internal/routes/register"
	"github.com/vvkh/social-network/internal/routes/user"
	"github.com/vvkh/social-network/internal/routes/users"
	"github.com/vvkh/social-network/internal/templates"
)

func (s *server) setupRoutes(templatesDir string) {
	templates := templates.New(templatesDir, "bootstrap").Add("base.gohtml")

	s.handler.Get("/", index.Handle())
	s.handler.Get("/login/", login.Handle(templates))
	s.handler.Get("/register/", register.Handle(templates))
	s.handler.Get("/friends/", friends.Handle(templates))
	s.handler.Route("/users/", func(r chi.Router) {
		r.Get("/", users.Handle(templates))
		r.Get("/{userID}/", user.Handle(templates))
	})

}
