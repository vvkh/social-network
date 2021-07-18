package server

import (
	"github.com/go-chi/chi/v5"

	"github.com/vvkh/social-network/internal/users"

	"github.com/vvkh/social-network/internal/routes/friends"
	"github.com/vvkh/social-network/internal/routes/index"
	"github.com/vvkh/social-network/internal/routes/login"
	"github.com/vvkh/social-network/internal/routes/profile"
	"github.com/vvkh/social-network/internal/routes/profiles"
	"github.com/vvkh/social-network/internal/routes/register"
	"github.com/vvkh/social-network/internal/templates"
)

func (s *server) setupRoutes(templatesDir string, usersUseCase users.UseCase) {
	templates := templates.New(templatesDir, "bootstrap").Add("base.gohtml")

	s.handler.Get("/", index.Handle())
	s.handler.Route("/login/", func(r chi.Router) {
		r.Get("/", login.HandleGet(templates))
		r.Post("/", login.HandlePost(usersUseCase))
	})
	s.handler.Route("/register/", func(r chi.Router) {
		r.Get("/", register.HandleGet(templates))
		r.Post("/", register.HandlePost(usersUseCase))
	})
	s.handler.Get("/friends/", friends.Handle(templates))
	s.handler.Route("/profiles/", func(r chi.Router) {
		r.Get("/", profiles.Handle(templates))
		r.Get("/{profileID}/", profile.Handle(templates))
	})

}
