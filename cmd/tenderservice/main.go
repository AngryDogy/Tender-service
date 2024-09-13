package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
	"tenderservice/handler"
	"tenderservice/logger"
	"tenderservice/repository"
	"tenderservice/server"
	"tenderservice/service"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logger.ErrorLogger.Fatal("failed to load env vars: ", err)
	}

	postgresRepository := repository.NewPostgresRepository(os.Getenv("POSTGRES_CONN"))
	err = postgresRepository.Connect()
	if err != nil {
		logger.ErrorLogger.Fatal("failed to connect to a database: ", err)
	}
	defer postgresRepository.Close()

	err = postgresRepository.Initialize("scheme.sql")
	if err != nil {
		logger.ErrorLogger.Fatal("failed to initialize a database: ", err)
	}

	authService := service.NewAuthService(postgresRepository)
	handlers := []handler.Handler{
		handler.NewTenderHandler(
			service.NewTenderService(postgresRepository),
			authService,
		),
		handler.NewBidHandler(
			service.NewBidService(postgresRepository),
			authService),
	}
	tenderServer := server.NewServer(
		gin.Default(),
		handlers,
	)

	tenderServer.Run(os.Getenv("SERVER_ADDRESS"))
}
