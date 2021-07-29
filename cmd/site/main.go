package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/vvkh/social-network/internal/config"
	"github.com/vvkh/social-network/internal/db"
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

	appConfig := config.NewFromEnv()

	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}

	sugarLogger := logger.Sugar()

	appDB, err := db.New(appConfig.DBUrl)
	if err != nil {
		return err
	}

	profilesRepo := profilesRepository.New(appDB)
	profilesUC := profilesUseCase.New(profilesRepo)

	usersRepo := usersRepository.New(appDB)
	usersUC := usersUseCase.New(profilesUC, usersRepo, appConfig.AuthSecret)

	friendshipRepo := friendshipRepository.New(appDB)
	friendshipUC := friendshipUseCase.New(friendshipRepo, profilesUC)

	s := server.New(sugarLogger, appConfig.ServerAddress, appConfig.TemplatesDir, usersUC, profilesUC, friendshipUC)
	return s.Start()
}
