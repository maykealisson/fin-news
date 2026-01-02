package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/maykealisson/fin-news/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockNoticiaService é um mock para INoticiaService
type MockNoticiaService struct {
	mock.Mock
}

func (m *MockNoticiaService) BuscarNoticias(ativo string) ([]services.Noticia, error) {
	args := m.Called(ativo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]services.Noticia), args.Error(1)
}

func TestBuscarNoticias_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockNoticiaService)
	controller := NewNoticiaController(mockService)

	// Setup Request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/noticias?ativo=PETR4", nil)

	// Mock Service Expectation
	expectedNews := []services.Noticia{
		{Titulo: "Test News", Link: "http://test.com"},
	}
	mockService.On("BuscarNoticias", "PETR4").Return(expectedNews, nil)

	// Execute
	controller.BuscarNoticias(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	// We could also check JSON body explicitly but status is good start
	mockService.AssertExpectations(t)
}

func TestBuscarNoticias_MissingParam(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockNoticiaService)
	controller := NewNoticiaController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/noticias", nil) // Missing active

	controller.BuscarNoticias(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "BuscarNoticias")
}

func TestBuscarNoticias_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockNoticiaService)
	controller := NewNoticiaController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/noticias?ativo=PETR4", nil)

	mockService.On("BuscarNoticias", "PETR4").Return(nil, errors.New("service failure"))

	controller.BuscarNoticias(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}
