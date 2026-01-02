package services

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/maykealisson/fin-news/clients"
	log "github.com/sirupsen/logrus"
)

// Constantes para configuração
const (
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

// INoticiaService define a interface para o serviço de notícias
type INoticiaService interface {
	BuscarNoticias(ativo string) ([]Noticia, error)
}

// NoticiaService implementa INoticiaService
type NoticiaService struct {
	client clients.IFinlightClient
	logger *log.Entry
}

// NewNoticiaService cria uma nova instância do serviço de notícias
func NewNoticiaService(client clients.IFinlightClient) *NoticiaService {
	return &NoticiaService{
		client: client,
		logger: log.WithField("service", "noticia"),
	}
}

func (s *NoticiaService) BuscarNoticias(ativo string) ([]Noticia, error) {
	logger := s.logger.WithField("ativo", ativo)

	// Implementa retry com backoff exponencial
	operation := func() ([]Noticia, error) {
		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()

		articles, err := s.client.BuscarArtigos(ctx, ativo)
		if err != nil {
			return nil, err
		}

		return s.converterParaNoticias(articles), nil
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

func (s *NoticiaService) converterParaNoticias(articles []clients.ArticleResponse) []Noticia {
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
