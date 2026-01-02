package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockHTTPClient é um mock para HTTPClient
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestBuscarArtigos_CacheHit(t *testing.T) {
	db, mockRedis := redismock.NewClientMock()
	mockHTTP := new(MockHTTPClient)
	apiKey := "test-key"
	client := NewFinlightClient(apiKey, db, mockHTTP)

	ctx := context.Background()
	query := "PETR4"
	expectedArticles := []ArticleResponse{
		{Title: "Noticia 1", Link: "http://link1.com"},
	}
	jsonArticles, _ := json.Marshal(expectedArticles)

	// Expect Redis Get to return data
	mockRedis.ExpectGet("noticias:PETR4").SetVal(string(jsonArticles))

	articles, err := client.BuscarArtigos(ctx, query)

	assert.NoError(t, err)
	assert.Equal(t, len(expectedArticles), len(articles))
	assert.Equal(t, expectedArticles[0].Title, articles[0].Title)

	// Ensure HTTP was NOT called
	mockHTTP.AssertNotCalled(t, "Do")
	assert.NoError(t, mockRedis.ExpectationsWereMet())
}

func TestBuscarArtigos_CacheMiss_APIHit(t *testing.T) {
	db, mockRedis := redismock.NewClientMock()
	mockHTTP := new(MockHTTPClient)
	apiKey := "test-key"
	client := NewFinlightClient(apiKey, db, mockHTTP)

	ctx := context.Background()
	query := "PETR4"

	// Cache miss
	mockRedis.ExpectGet("noticias:PETR4").RedisNil()

	// API Response
	apiResponse := ArticlesResponse{
		Articles: []ArticleResponse{
			{Title: "Noticia API", Link: "http://api.com"},
		},
	}
	jsonResponse, _ := json.Marshal(apiResponse)

	// Mock HTTP Success
	responseBody := io.NopCloser(bytes.NewReader(jsonResponse))
	mockHTTP.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: 200,
		Body:       responseBody,
	}, nil)

	// Expect Redis Set
	// We expect the JSON content of expectation (articles from API)
	expectedJSON, _ := json.Marshal(apiResponse.Articles)
	mockRedis.ExpectSet("noticias:PETR4", expectedJSON, 24*time.Hour).SetVal("OK")

	articles, err := client.BuscarArtigos(ctx, query)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(articles))
	assert.Equal(t, "Noticia API", articles[0].Title)

	mockHTTP.AssertExpectations(t)
	assert.NoError(t, mockRedis.ExpectationsWereMet())
}

func TestBuscarArtigos_APIFailure(t *testing.T) {
	db, mockRedis := redismock.NewClientMock()
	mockHTTP := new(MockHTTPClient)
	apiKey := "test-key"
	client := NewFinlightClient(apiKey, db, mockHTTP)

	ctx := context.Background()
	query := "PETR4"

	// Cache miss
	mockRedis.ExpectGet("noticias:PETR4").RedisNil()

	// Mock HTTP Error
	mockHTTP.On("Do", mock.Anything).Return(nil, errors.New("network error"))

	articles, err := client.BuscarArtigos(ctx, query)

	assert.Error(t, err)
	assert.Nil(t, articles)

	mockHTTP.AssertExpectations(t)
}
