package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	defaultTimeout = 10 * time.Second
	baseURL        = "https://api.finlight.me/v2/articles"
)

type FinlightClient struct {
	client *http.Client
	apiKey string
	logger *log.Entry
}

type ArticleRequest struct {
	Query    string `json:"query"`
	PageSize string `json:"pageSize"`
}

type ArticleResponse struct {
	Link        string   `json:"link"`
	Title       string   `json:"title"`
	Summary     string   `json:"summary"`
	PublishDate string   `json:"publishDate"`
	Images      []string `json:"images"`
}

type ArticlesResponse struct {
	Articles []ArticleResponse `json:"articles"`
}

func NewFinlightClient(apiKey string) *FinlightClient {
	return &FinlightClient{
		client: &http.Client{
			Timeout: defaultTimeout,
		},
		apiKey: apiKey,
		logger: log.WithField("client", "finlight"),
	}
}

func (f *FinlightClient) BuscarArtigos(ctx context.Context, query string) ([]ArticleResponse, error) {
	reqBody := ArticleRequest{
		Query:    query,
		PageSize: "5",
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter body: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", baseURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", f.apiKey)

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro na requisição: %v", err)
	}
	defer resp.Body.Close()

	f.logger.WithField("status", resp.StatusCode).Debug("Response received")

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("erro na API (status %d): %s", resp.StatusCode, string(body))
	}

	var response ArticlesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %v", err)
	}

	return response.Articles, nil
}
