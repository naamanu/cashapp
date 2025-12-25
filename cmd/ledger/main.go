package main

import (
	"cashapp/core"
	"cashapp/core/database"
	"cashapp/internal/ledger/api"
	"cashapp/internal/ledger/models"
	"cashapp/internal/ledger/repository"
	"cashapp/internal/ledger/service"

	"go.uber.org/zap"

	_ "cashapp/docs"
)

// @title CashApp Ledger Service
// @version 1.0
// @description Payment processing service
// @BasePath /
func main() {
	config := core.NewConfig()
	// Override port if needed, or use Env. For now sharing same config logic.
	// Ideally we set PORT env var differently for each service.
	core.InitLogger(config.ENVIRONMENT)

	pg, err := database.NewPostgres(config)
	if err != nil {
		core.Log.Fatal("failed to initialize postgres database", zap.Error(err))
	}

	err = database.RunMigrations(pg, &models.Transaction{}, &models.TransactionEvent{}, &models.PaymentRequest{})
	if err != nil {
		core.Log.Fatal("failed to run migrations", zap.Error(err))
	}

	repo := repository.New(pg)
	svc := service.New(repo, config)
	server := core.NewHTTPServer(config)

	api.RegisterPaymentRoutes(server.Engine, svc)
	server.Start()
}
