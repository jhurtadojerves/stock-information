package ports

import "stock-information/internal/domain"

type StockRepository interface {
	SaveRecommendations([]domain.Recommendation) error
	GetAllRecommendations() ([]domain.Recommendation, error)
	GetRecommendations(page int, limit int, sort string) ([]domain.Recommendation, int, error)
	GetRecommendationsByCompany(company string, page int, limit int, sort string) ([]domain.Recommendation, int, error)
	GetRecommendationByTicker(ticker string) (domain.Recommendation, error)
}
