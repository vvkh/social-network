package server

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/vvkh/social-network/internal/domain/users"
	"github.com/vvkh/social-network/internal/middlewares"
	"github.com/vvkh/social-network/internal/permissions"
	"github.com/vvkh/social-network/internal/routes/friends"
	"github.com/vvkh/social-network/internal/routes/index"
	"github.com/vvkh/social-network/internal/routes/login"
	"github.com/vvkh/social-network/internal/routes/logout"
	"github.com/vvkh/social-network/internal/routes/profile"
	"github.com/vvkh/social-network/internal/routes/profiles"
	"github.com/vvkh/social-network/internal/routes/register"
	"github.com/vvkh/social-network/internal/templates"
)

const (
	defaultHandlerTimeout = 60 * time.Second
)

func (s *server) setupRoutes(templatesDir string, usersUseCase users.UseCase) {
	templates := templates.New(templatesDir, "bootstrap").Add("base.gohtml")

	s.handler.Use(middleware.RequestID)
	s.handler.Use(middleware.RealIP)
	s.handler.Use(middleware.Logger)
	s.handler.Use(middleware.Recoverer)
	s.handler.Use(middleware.Timeout(defaultHandlerTimeout))
	s.handler.Use(middlewares.AuthenticateUser(usersUseCase))

	authRequired := permissions.AuthRequired("/login/")

	s.handler.Get("/", authRequired(index.Handle()))
	s.handler.Route("/login/", func(r chi.Router) {
		r.Get("/", login.HandleGet(templates))
		r.Post("/", login.HandlePost(usersUseCase, "/"))
	})
	s.handler.Get("/logout/", logout.HandleGet("/"))
	s.handler.Route("/register/", func(r chi.Router) {
		r.Get("/", register.HandleGet(templates))
		r.Post("/", register.HandlePost(usersUseCase, "/login/"))
	})
	s.handler.Get("/friends/", authRequired(friends.Handle(templates)))
	s.handler.Route("/profiles/", func(r chi.Router) {
		r.Get("/", authRequired(profiles.Handle(templates)))
		r.Get("/{profileID}/", authRequired(profile.Handle(templates)))
	})

}
