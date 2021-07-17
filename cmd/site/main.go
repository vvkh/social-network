package main

import (
	"fmt"

	"github.com/joho/godotenv"

	"github.com/vvkh/social-network/internal/server"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("error while loading dotenv: %v", err)
		return
	}

	s, err := server.NewDefault()
	if err != nil {
		fmt.Printf("error while initialising server: %v", err)
		return
	}

	if err := s.Start(); err != nil {
		fmt.Printf("server returned error: %v", err)
	}
}
