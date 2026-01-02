package services

import (
	"context"
	"errors"
	"testing"

	"github.com/maykealisson/fin-news/clients"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFinlightClient é um mock para IFinlightClient
type MockFinlightClient struct {
	mock.Mock
}

func (m *MockFinlightClient) BuscarArtigos(ctx context.Context, query string) ([]clients.ArticleResponse, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]clients.ArticleResponse), args.Error(1)
}

func TestBuscarNoticias_Success(t *testing.T) {
	mockClient := new(MockFinlightClient)
	service := NewNoticiaService(mockClient)
	ativo := "PETR4"

	expectedArticles := []clients.ArticleResponse{
		{Title: "Title 1", Summary: "Summary 1", Link: "Link 1", PublishDate: "2023-01-01"},
	}

	mockClient.On("BuscarArtigos", mock.Anything, ativo).Return(expectedArticles, nil)

	noticias, err := service.BuscarNoticias(ativo)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(noticias))
	assert.Equal(t, "Title 1", noticias[0].Titulo)

	mockClient.AssertExpectations(t)
}

func TestBuscarNoticias_Error(t *testing.T) {
	mockClient := new(MockFinlightClient)
	service := NewNoticiaService(mockClient)
	ativo := "PETR4"

	// Mock error - backoff will retry, so we might see multiple calls if we don't mock carefully using .Times() or just allowing anything
	// However, backoff usually retries on error. We want to ensure failure eventually returns error
	// The service uses backoff.Retry. If mock returns error every time, it should fail.

	mockClient.On("BuscarArtigos", mock.Anything, ativo).Return(nil, errors.New("api error"))

	noticias, err := service.BuscarNoticias(ativo)

	assert.Error(t, err)
	assert.Nil(t, noticias)

	// Ensure it was called at least once (likely multiple times due to retry)
	mockClient.AssertCalled(t, "BuscarArtigos", mock.Anything, ativo)
}
