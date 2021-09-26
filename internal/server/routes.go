package server

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/vvkh/social-network/internal/domain/friendship"
	profilesDomain "github.com/vvkh/social-network/internal/domain/profiles"
	"github.com/vvkh/social-network/internal/domain/users"
	"github.com/vvkh/social-network/internal/middlewares"
	navbar "github.com/vvkh/social-network/internal/navbar"
	"github.com/vvkh/social-network/internal/permissions"
	"github.com/vvkh/social-network/internal/routes/friend"
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

func (s *server) setupRoutes(log *zap.SugaredLogger, templatesDir string, usersUseCase users.UseCase, profilesUseCase profilesDomain.UseCase, friendshipUseCase friendship.UseCase) {
	navbar := navbar.New(friendshipUseCase, nil)
	templates := templates.New(templatesDir, "bootstrap").Add("base.gohtml")

	s.handler.Use(middleware.RequestID)
	s.handler.Use(middleware.RealIP)
	s.handler.Use(middleware.Logger)
	s.handler.Use(middleware.Recoverer)
	s.handler.Use(middleware.Timeout(defaultHandlerTimeout))
	s.handler.Use(middlewares.AuthenticateUser(log, usersUseCase, profilesUseCase))

	authRequired := permissions.AuthRequired("/login/")

	s.handler.Get("/", authRequired(index.Handle()))
	s.handler.Route("/login/", func(r chi.Router) {
		r.Get("/", login.HandleGet(templates))
		r.Post("/", login.HandlePost(log, usersUseCase, "/", templates))
	})
	s.handler.Get("/logout/", logout.HandleGet("/"))
	s.handler.Route("/register/", func(r chi.Router) {
		r.Get("/", register.HandleGet(templates))
		r.Post("/", register.HandlePost(log, usersUseCase, "/login/", templates))
	})
	s.handler.Route("/friends/", func(r chi.Router) {
		r.Get("/", authRequired(friends.Handle(friendshipUseCase, navbar, templates)))
		r.Post("/{profile:[0-9]+}/stop/", friend.HandleStop(friendshipUseCase))
		r.Route("/requests", func(r chi.Router) {
			r.Get("/", authRequired(friends_requests.Handle(friendshipUseCase, navbar, templates)))
			r.Post("/{profileFrom:[0-9]+}/accept/", authRequired(friends_requests.HandlePostAccept(friendshipUseCase, "/friends/requests/")))
			r.Post("/{profileFrom:[0-9]+}/decline/", authRequired(friends_requests.HandlePostDecline(friendshipUseCase, "/friends/requests/")))
			r.Post("/{profileTo:[0-9]+}/create/", authRequired(friends_requests.HandleCreate(friendshipUseCase)))
		})

	})
	s.handler.Route("/profiles/", func(r chi.Router) {
		r.Get("/", authRequired(profiles.Handle(log, profilesUseCase, navbar, templates)))
		r.Get("/{profileID:[0-9]+}/", authRequired(profile.Handle(profilesUseCase, friendshipUseCase, navbar, templates)))
	})
}
