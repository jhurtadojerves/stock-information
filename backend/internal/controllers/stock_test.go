package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"stock-information/internal/controllers"
	"stock-information/internal/services"
	"stock-information/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetAllRecommendations(t *testing.T) {
	stockController := &controllers.StockController{
		StockService: &services.StockService{
			Repo: &mocks.MockStockRepository{},
		},
	}

	r := gin.Default()
	r.GET("/api/v1/stocks", stockController.GetAllRecommendations)

	req, _ := http.NewRequest("GET", "/api/v1/stocks?page=1&limit=10&sort=asc", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	assert.Contains(t, w.Body.String(), "AAPL")
	assert.Contains(t, w.Body.String(), "Apple")
	assert.Contains(t, w.Body.String(), "GOOG")
	assert.Contains(t, w.Body.String(), "Google")
}

func TestGetBestRecommendation(t *testing.T) {
	stockController := &controllers.StockController{
		StockService: &services.StockService{
			Repo: &mocks.MockStockRepository{},
		},
	}

	r := gin.Default()
	r.GET("/api/v1/stocks/recommendation", stockController.GetBestRecommendation)

	req, _ := http.NewRequest("GET", "/api/v1/stocks/recommendation", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	assert.Contains(t, w.Body.String(), "GOOG")
	assert.Contains(t, w.Body.String(), "Google")
}

func TestGetRecommendationByTicker(t *testing.T) {
	stockController := &controllers.StockController{
		StockService: &services.StockService{
			Repo: &mocks.MockStockRepository{},
		},
	}

	r := gin.Default()
	r.GET("/api/v1/stocks/:ticker", stockController.GetRecommendationByTicker)

	req, _ := http.NewRequest("GET", "/api/v1/stocks/AAPL", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	assert.Contains(t, w.Body.String(), "AAPL")
	assert.Contains(t, w.Body.String(), "Apple")
}
