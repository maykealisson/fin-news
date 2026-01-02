package routes

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/maykealisson/fin-news/clients"
	"github.com/maykealisson/fin-news/config"
	"github.com/maykealisson/fin-news/controllers"
	"github.com/maykealisson/fin-news/services"
	log "github.com/sirupsen/logrus"
)

func HandlerRequests() {
	// Setup Logging
	log.SetFormatter(&log.JSONFormatter{})

	// Load env
	if err := godotenv.Load(); err != nil {
		log.Warn("Arquivo .env não encontrado, usando variáveis de ambiente")
	}

	// Dependencies
	apiKey := os.Getenv("FINLIGHT_KEY")
	if apiKey == "" {
		log.Fatal("FINLIGHT_KEY não encontrada")
	}

	// Initialize Redis
	redisClient := config.NewRedisClient()
	defer redisClient.Close()

	// Initialize Layers
	finlightClient := clients.NewFinlightClient(apiKey, redisClient, nil) // nil uses default http.Client
	noticiaService := services.NewNoticiaService(finlightClient)
	noticiaController := controllers.NewNoticiaController(noticiaService)

	server := config.SetupGin()

	server.GET("/fin-news/v1/noticias", noticiaController.BuscarNoticias)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	log.Infof("Server starting on port %s", port)
	server.Run(":" + port)
}
