package main

import (
	"github.com/gin-gonic/gin"
	"tenderservice/server"
)

func main() {
	/*err := godotenv.Load()
	if err != nil {
		logger.ErrorLogger.Fatal("failed to load env vars: ", err)
	}*/

	/*postgresRepository := repository.NewPostgresRepository(os.Getenv("POSTGRES_CONN"))
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
	}*/
	tenderServer := server.NewServer(
		gin.Default(),
		handlers,
	)

	tenderServer.Run("0.0.0.0:8080")
}
