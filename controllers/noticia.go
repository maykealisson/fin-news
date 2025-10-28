package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/maykealisson/fin-news/services"
	log "github.com/sirupsen/logrus"
)

// NoticiaResponse representa a estrutura da resposta
type NoticiaResponse struct {
	Noticias []services.Noticia `json:"noticias"`
}

// ErrorResponse representa a estrutura de erro
type ErrorResponse struct {
	Error string `json:"error"`
}

// BuscarNoticias godoc
// @Summary Busca notícias relacionadas a um ativo
// @Description Retorna uma lista de notícias financeiras relacionadas ao ativo informado
// @Tags noticias
// @Accept json
// @Produce json
// @Param ativo query string true "Código do ativo (ex: PETR4)"
// @Success 200 {object} NoticiaResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /noticias [get]
func BuscarNoticias(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"handler": "BuscarNoticias",
		"method":  c.Request.Method,
		"path":    c.Request.URL.Path,
	})

	ativo := strings.TrimSpace(c.Query("ativo"))
	if ativo == "" {
		logger.Warn("Parâmetro 'ativo' não informado")
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Parâmetro 'ativo' é obrigatório",
		})
		return
	}

	logger = logger.WithField("ativo", ativo)
	logger.Info("Buscando notícias")

	news, err := services.BuscarNoticias(ativo)
	if err != nil {
		logger.WithError(err).Error("Erro ao buscar notícias")
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Erro interno ao processar a requisição",
		})
		return
	}

	logger.WithField("quantidade", len(news)).Info("Notícias encontradas")
	c.JSON(http.StatusOK, NoticiaResponse{
		Noticias: news,
	})
}
