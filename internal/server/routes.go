package server

import (
	"github.com/go-chi/chi/v5"

	"github.com/vvkh/social-network/internal/routes/friends"
	"github.com/vvkh/social-network/internal/routes/index"
	"github.com/vvkh/social-network/internal/routes/login"
	"github.com/vvkh/social-network/internal/routes/profile"
	"github.com/vvkh/social-network/internal/routes/profiles"
	"github.com/vvkh/social-network/internal/routes/register"
	"github.com/vvkh/social-network/internal/templates"
)

func (s *server) setupRoutes(templatesDir string) {
	templates := templates.New(templatesDir, "bootstrap").Add("base.gohtml")

	s.handler.Get("/", index.Handle())
	s.handler.Get("/login/", login.Handle(templates))
	s.handler.Get("/register/", register.Handle(templates))
	s.handler.Get("/friends/", friends.Handle(templates))
	s.handler.Route("/profiles/", func(r chi.Router) {
		r.Get("/", profiles.Handle(templates))
		r.Get("/{profileID}/", profile.Handle(templates))
	})

}
