package mocks

import (
	"fmt"
	"stock-information/internal/domain"
)

type MockStockRepository struct{}

func (m *MockStockRepository) SaveRecommendations(recs []domain.Recommendation) error {
	return nil
}

func (m *MockStockRepository) GetRecommendations(page int, limit int, sort string) ([]domain.Recommendation, int, error) {
	recs := []domain.Recommendation{
		{Ticker: "AAPL", Company: "Apple", Action: "Buy", RatingFrom: "Neutral", RatingTo: "Buy", TargetFrom: 100.0, TargetTo: 150.0},
		{Ticker: "GOOG", Company: "Google", Action: "Hold", RatingFrom: "Neutral", RatingTo: "Hold", TargetFrom: 2000.0, TargetTo: 2500.0},
	}
	return recs, len(recs), nil
}

func (m *MockStockRepository) GetRecommendationsByCompany(company string, page int, limit int, sort string) ([]domain.Recommendation, int, error) {
	if company == "Apple" {
		recs := []domain.Recommendation{
			{Ticker: "AAPL", Company: "Apple", Action: "Buy", RatingFrom: "Neutral", RatingTo: "Buy", TargetFrom: 100.0, TargetTo: 150.0},
		}
		return recs, len(recs), nil
	}
	return nil, 0, fmt.Errorf("company not found")
}

func (m *MockStockRepository) GetRecommendationByTicker(ticker string) (domain.Recommendation, error) {
	if ticker == "AAPL" {
		return domain.Recommendation{Ticker: "AAPL", Company: "Apple", Action: "Buy", RatingFrom: "Neutral", RatingTo: "Buy", TargetFrom: 100.0, TargetTo: 150.0}, nil
	}
	return domain.Recommendation{}, fmt.Errorf("recommendation not found for ticker %s", ticker)
}

func (m *MockStockRepository) GetAllRecommendations() ([]domain.Recommendation, error) {
	return []domain.Recommendation{
		{Ticker: "AAPL", Company: "Apple", Action: "Buy", RatingFrom: "Neutral", RatingTo: "Buy", TargetFrom: 100.0, TargetTo: 150.0},
		{Ticker: "GOOG", Company: "Google", Action: "Hold", RatingFrom: "Neutral", RatingTo: "Hold", TargetFrom: 2000.0, TargetTo: 2500.0},
	}, nil
}
