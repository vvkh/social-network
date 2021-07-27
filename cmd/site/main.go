package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	friendshipRepository "github.com/vvkh/social-network/internal/domain/friendship/repository"
	friendshipUseCase "github.com/vvkh/social-network/internal/domain/friendship/usecase"
	profilesRepository "github.com/vvkh/social-network/internal/domain/profiles/repository"
	profilesUseCase "github.com/vvkh/social-network/internal/domain/profiles/usecase"
	usersRepository "github.com/vvkh/social-network/internal/domain/users/repository"
	usersUseCase "github.com/vvkh/social-network/internal/domain/users/usecase"
	"github.com/vvkh/social-network/internal/heroku"
	"github.com/vvkh/social-network/internal/server"
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

	err = heroku.ConvertEnv()
	if err != nil {
		return err
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}

	sugarLogger := logger.Sugar()

	// TODO: single db instance
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

	friendshipRepo, err := friendshipRepository.NewDefault()
	if err != nil {
		return err
	}

	friendshipUC := friendshipUseCase.New(friendshipRepo, profilesUC)
	s, err := server.NewFromEnv(sugarLogger, usersUC, profilesUC, friendshipUC)
	if err != nil {
		return err
	}

	return s.Start()
}
