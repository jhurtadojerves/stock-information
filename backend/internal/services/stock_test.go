package services_test

import (
	"stock-information/internal/domain"
	"stock-information/internal/services"
	"stock-information/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRecommendations(t *testing.T) {
	mockRepo := &mocks.MockStockRepository{}
	stockService := &services.StockService{Repo: mockRepo}
	recs, totalCount, err := stockService.GetRecommendations(1, 10, "asc")

	assert.NoError(t, err)
	assert.Equal(t, 2, len(recs))
	assert.Equal(t, 2, totalCount)
}

func TestGetRecommendationsByCompany(t *testing.T) {
	mockRepo := &mocks.MockStockRepository{}
	stockService := &services.StockService{Repo: mockRepo}

	recs, totalCount, err := stockService.GetRecommendationsByCompany("Apple", 1, 10, "asc")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(recs))
	assert.Equal(t, 1, totalCount)

	recs, _, err = stockService.GetRecommendationsByCompany("NonExistent", 1, 10, "asc")
	assert.Error(t, err)
	assert.Equal(t, 0, len(recs))
}

func TestGetRecommendationByTicker(t *testing.T) {
	mockRepo := &mocks.MockStockRepository{}
	stockService := &services.StockService{Repo: mockRepo}

	rec, err := stockService.GetRecommendationByTicker("AAPL")
	assert.NoError(t, err)
	assert.Equal(t, "AAPL", rec.Ticker)
	assert.Equal(t, "Apple", rec.Company)

	rec, err = stockService.GetRecommendationByTicker("NONEXISTENT")
	assert.Error(t, err)
	assert.Equal(t, domain.Recommendation{}, rec)
}
func TestGetBestRecommendation(t *testing.T) {
	mockRepo := &mocks.MockStockRepository{}
	stockService := &services.StockService{Repo: mockRepo}

	bestRecommendation, err := stockService.GetBestRecommendation()

	assert.NoError(t, err)
	assert.Equal(t, "GOOG", bestRecommendation.Ticker)
}
