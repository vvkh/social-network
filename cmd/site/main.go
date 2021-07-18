package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	profilesRepository "github.com/vvkh/social-network/internal/profiles/repository"
	profilesUseCase "github.com/vvkh/social-network/internal/profiles/usecase"
	"github.com/vvkh/social-network/internal/server"
	usersRepository "github.com/vvkh/social-network/internal/users/repository"
	usersUseCase "github.com/vvkh/social-network/internal/users/usecase"
)

func main() {
	err := run()
	if err != nil {
		fmt.Printf("server returned err: %v", err)
		os.Exit(1)
	}
}

func run() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	profilesRepo, err := profilesRepository.NewDefault()
	if err != nil {
		return err
	}

	profilesUC := profilesUseCase.New(profilesRepo)
	usersRepo, err := usersRepository.NewDefault()
	if err != nil {
		return err
	}

	usersUC := usersUseCase.NewFromEnv(profilesUC, usersRepo)
	s, err := server.NewFromEnv(usersUC)
	if err != nil {
		return err
	}

	return s.Start()
}
