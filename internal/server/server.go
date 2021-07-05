package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type server struct {
	handler *http.ServeMux
	address string
}

func NewDefault() (*server, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error while parsing env: %w", err)
	}

	address := os.Getenv("SERVER_ADDRESS")
	s := New(address)
	return s, nil
}

func New(address string) *server {
	handler := http.NewServeMux()
	s := server{
		handler: handler,
		address: address,
	}
	s.setupRoutes()
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
