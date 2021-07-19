package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"github.com/vvkh/social-network/internal/domain/users"
)

type server struct {
	handler *chi.Mux
	address string
}

func NewFromEnv(usersUseCase users.UseCase) (*server, error) {
	address := os.Getenv("SERVER_ADDRESS")
	templatesDir := os.Getenv("TEMPLATES_DIR")
	s := New(address, templatesDir, usersUseCase)
	return s, nil
}

func New(address string, tempalatesDir string, usersUseCase users.UseCase) *server {
	router := chi.NewRouter()
	s := server{
		handler: router,
		address: address,
	}
	s.setupRoutes(tempalatesDir, usersUseCase)
	return &s
}

func (s *server) Start() error {
	httpServer := http.Server{
		Handler: s.handler,
		Addr:    s.address,
	}
	fmt.Printf("starting server on %s", s.address)
	return httpServer.ListenAndServe()
}

func (s *server) Handle(writer http.ResponseWriter, request *http.Request) {
	s.handler.ServeHTTP(writer, request)
}
