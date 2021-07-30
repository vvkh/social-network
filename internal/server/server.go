package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/vvkh/social-network/internal/domain/friendship"
	"github.com/vvkh/social-network/internal/domain/profiles"
	"github.com/vvkh/social-network/internal/domain/users"
)

type server struct {
	handler *chi.Mux
	address string
	log     *zap.SugaredLogger
}

func New(log *zap.SugaredLogger, address string, tempalatesDir string, usersUseCase users.UseCase, profilesUseCase profiles.UseCase, friendshipUseCase friendship.UseCase) *server {
	router := chi.NewRouter()
	s := server{
		handler: router,
		address: address,
		log:     log,
	}
	s.setupRoutes(log, tempalatesDir, usersUseCase, profilesUseCase, friendshipUseCase)
	return &s
}

func (s *server) Start() error {
	s.log.Infow("starting server", "address", s.address)

	httpServer := http.Server{
		Handler: s.handler,
		Addr:    s.address,
	}
	return httpServer.ListenAndServe()
}

func (s *server) Handle(writer http.ResponseWriter, request *http.Request) {
	s.handler.ServeHTTP(writer, request)
}
