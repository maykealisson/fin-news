package config

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func SetupGin() *gin.Engine {
	logger := log.WithFields(log.Fields{
		"service": "SetupGin",
	})
	if err := godotenv.Load(); err != nil {
		logger.Error("Erro ao carregar .env")
		os.Exit(1)
	}

	environment := os.Getenv("ENV")
	if environment == "" {
		logger.Error("ENV não encontrada no .env")
		os.Exit(1)
	}
	if environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Create a new Gin instance without default middleware
	r := gin.New()

	// Add necessary middleware manually
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// // Set trusted proxies
	r.SetTrustedProxies([]string{"127.0.0.1"})

	return r
}
