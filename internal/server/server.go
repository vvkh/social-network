package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"github.com/joho/godotenv"
)

type server struct {
	handler *chi.Mux
	address string
}

func NewDefault() (*server, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error while parsing env: %w", err)
	}

	address := os.Getenv("SERVER_ADDRESS")
	templatesDir := os.Getenv("TEMPLATES_DIR")
	s := New(address, templatesDir)
	return s, nil
}

func New(address string, tempalatesDir string) *server {
	router := chi.NewRouter()
	s := server{
		handler: router,
		address: address,
	}
	s.setupRoutes(tempalatesDir)
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
