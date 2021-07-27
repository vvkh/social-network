package server

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/vvkh/social-network/internal/domain/friendship"
	"github.com/vvkh/social-network/internal/domain/profiles"
	"github.com/vvkh/social-network/internal/domain/users"
)

type server struct {
	handler *chi.Mux
	address string
}

func NewFromEnv(log *zap.SugaredLogger, usersUseCase users.UseCase, profilesUseCase profiles.UseCase, friendshipUseCase friendship.UseCase) (*server, error) {
	address := os.Getenv("SERVER_ADDRESS")
	templatesDir := os.Getenv("TEMPLATES_DIR")
	log.Infow("starting server", "address", address, "templateDir", templatesDir)
	s := New(log, address, templatesDir, usersUseCase, profilesUseCase, friendshipUseCase)
	return s, nil
}

func New(log *zap.SugaredLogger, address string, tempalatesDir string, usersUseCase users.UseCase, profilesUseCase profiles.UseCase, friendshipUseCase friendship.UseCase) *server {
	router := chi.NewRouter()
	s := server{
		handler: router,
		address: address,
	}
	s.setupRoutes(log, tempalatesDir, usersUseCase, profilesUseCase, friendshipUseCase)
	return &s
}

func (s *server) Start() error {
	httpServer := http.Server{
		Handler: s.handler,
		Addr:    s.address,
	}
	return httpServer.ListenAndServe()
}

func (s *server) Handle(writer http.ResponseWriter, request *http.Request) {
	s.handler.ServeHTTP(writer, request)
}
