package server

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/vvkh/social-network/internal/domain/friendship"
	profilesDomain "github.com/vvkh/social-network/internal/domain/profiles"
	"github.com/vvkh/social-network/internal/domain/users"
	"github.com/vvkh/social-network/internal/middlewares"
	"github.com/vvkh/social-network/internal/permissions"
	"github.com/vvkh/social-network/internal/routes/friends"
	"github.com/vvkh/social-network/internal/routes/friends_requests"
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

func (s *server) setupRoutes(templatesDir string, usersUseCase users.UseCase, profilesUseCase profilesDomain.UseCase, friendshipUseCase friendship.UseCase) {
	templates := templates.New(templatesDir, "bootstrap").Add("base.gohtml")

	s.handler.Use(middleware.RequestID)
	s.handler.Use(middleware.RealIP)
	s.handler.Use(middleware.Logger)
	s.handler.Use(middleware.Recoverer)
	s.handler.Use(middleware.Timeout(defaultHandlerTimeout))
	s.handler.Use(middlewares.AuthenticateUser(usersUseCase, profilesUseCase))

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
	s.handler.Route("/friends/", func(r chi.Router) {
		r.Get("/", authRequired(friends.Handle(friendshipUseCase, templates)))
		r.Route("/requests", func(r chi.Router) {
			r.Get("/", authRequired(friends_requests.Handle(friendshipUseCase, templates)))
			r.Post("/{profileFrom:[0-9]+}/accept/", authRequired(friends_requests.HandlePostAccept(friendshipUseCase, "/friends/requests/")))
			r.Post("/{profileFrom:[0-9]+}/decline/", authRequired(friends_requests.HandlePostDecline(friendshipUseCase, "/friends/requests/")))
		})

	})
	s.handler.Route("/profiles/", func(r chi.Router) {
		r.Get("/", authRequired(profiles.Handle(profilesUseCase, templates)))
		r.Get("/{profileID:[0-9]+}/", authRequired(profile.Handle(profilesUseCase, templates)))
		r.Post("/{profileID:[0-9]+}/friendship/", authRequired(profile.HandlePost(friendshipUseCase)))
	})
}
