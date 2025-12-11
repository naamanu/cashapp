package main

import (
	"cashapp/core"
	"cashapp/core/database"
	"cashapp/internal/user/api"
	"cashapp/internal/user/models"
	"cashapp/internal/user/repository"
	"cashapp/internal/user/service"

	"go.uber.org/zap"

	_ "cashapp/docs"
)

// @title CashApp User Service
// @version 1.0
// @description User management service
// @BasePath /
func main() {
	config := core.NewConfig()
	core.InitLogger(config.ENVIRONMENT)

	pg, err := database.NewPostgres(config)
	if err != nil {
		core.Log.Fatal("failed to initialize postgres database", zap.Error(err))
	}

	err = database.RunMigrations(pg, &models.User{}, &models.Wallet{})
	if err != nil {
		core.Log.Fatal("failed to run migrations", zap.Error(err))
	}

	if config.RUN_SEEDS {
		models.RunSeeds(pg)
	}

	repo := repository.New(pg)
	svc := service.New(repo, config)
	server := core.NewHTTPServer(config)

	api.RegisterUserRoutes(server.Engine, svc)
	server.Start()
}
