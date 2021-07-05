package main

import (
	"fmt"

	"github.com/vvkh/social-network/internal/server"
)

func main() {
	s, err := server.NewDefault()
	if err != nil {
		fmt.Printf("error while initialising server: %v", err)
	}

	if err := s.Start(); err != nil {
		fmt.Printf("server returned error: %v", err)
	}
}
