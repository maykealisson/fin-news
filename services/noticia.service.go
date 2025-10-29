package services

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/joho/godotenv"
	"github.com/maykealisson/fin-news/clients"
	"github.com/maykealisson/fin-news/config"
	log "github.com/sirupsen/logrus"
)

// Constantes para configuração
const (
	maxRetries     = 3
	defaultTimeout = 10 * time.Second
)

// Noticia representa uma notícia financeira
type Noticia struct {
	Link   string   `json:"link"`
	Titulo string   `json:"titulo"`
	Resumo string   `json:"resumo"`
	Data   string   `json:"data"`
	Images []string `json:"images"`
}

func BuscarNoticias(ativo string) ([]Noticia, error) {
	logger := log.WithFields(log.Fields{
		"service": "noticia",
		"ativo":   ativo,
	})

	if err := godotenv.Load(); err != nil {
		logger.Error("Erro ao carregar .env")
		return nil, fmt.Errorf("erro ao carregar .env: %v", err)
	}

	apiKey := os.Getenv("FINLIGHT_KEY")
	if apiKey == "" {
		logger.Error("FINLIGHT_KEY não encontrada")
		return nil, fmt.Errorf("FINLIGHT_KEY não encontrada no .env")
	}

	redisClient := config.NewRedisClient()
	finlightClient := clients.NewFinlightClient(apiKey, redisClient)

	// Implementa retry com backoff exponencial
	operation := func() ([]Noticia, error) {
		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()

		articles, err := finlightClient.BuscarArtigos(ctx, ativo)
		if err != nil {
			return nil, err
		}

		return converterParaNoticias(articles), nil
	}

	// Configuração do backoff
	exponentialBackOff := backoff.NewExponentialBackOff()
	exponentialBackOff.MaxElapsedTime = 30 * time.Second

	var noticias []Noticia
	err := backoff.Retry(func() error {
		var err error
		noticias, err = operation()
		return err
	}, exponentialBackOff)

	if err != nil {
		logger.WithError(err).Error("Falha após todas as tentativas")
		return nil, err
	}

	logger.WithField("quantidade", len(noticias)).Info("Notícias recuperadas com sucesso")
	return noticias, nil
}

func converterParaNoticias(articles []clients.ArticleResponse) []Noticia {
	noticias := make([]Noticia, len(articles))
	for i, article := range articles {
		noticias[i] = Noticia{
			Link:   article.Link,
			Titulo: article.Title,
			Resumo: article.Summary,
			Data:   article.PublishDate,
			Images: article.Images,
		}
	}
	return noticias
}
