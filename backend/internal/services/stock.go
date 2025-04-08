package services

import (
	"sort"
	"stock-information/internal/domain"
	"stock-information/internal/ports"
)

type StockService struct {
	Repo ports.StockRepository
}

type scoredRecommendation struct {
	Recommendation domain.Recommendation
	Score          float64
}

func (s *StockService) SyncRecommendations(data []domain.Recommendation) error {
	return s.Repo.SaveRecommendations(data)
}

func (s *StockService) GetRecommendations(page int, limit int, sort string) ([]domain.Recommendation, int, error) {
	return s.Repo.GetRecommendations(page, limit, sort)
}

func (s *StockService) GetRecommendationsByCompany(company string, page int, limit int, sort string) ([]domain.Recommendation, int, error) {
	return s.Repo.GetRecommendationsByCompany(company, page, limit, sort)
}

func (s *StockService) GetRecommendationByTicker(ticker string) (domain.Recommendation, error) {
	return s.Repo.GetRecommendationByTicker(ticker)
}

func (s *StockService) GetBestRecommendation() (domain.Recommendation, error) {
	recommendations, err := s.Repo.GetAllRecommendations()

	if err != nil {
		return domain.Recommendation{}, err
	}

	var scoredRecommendations []scoredRecommendation

	for _, rec := range recommendations {
		score := calculateStockScore(rec)
		scoredRecommendations = append(scoredRecommendations, scoredRecommendation{
			Recommendation: rec,
			Score:          score,
		})
	}

	sort.Slice(scoredRecommendations, func(i, j int) bool {
		return scoredRecommendations[i].Score > scoredRecommendations[j].Score
	})

	return scoredRecommendations[0].Recommendation, nil

}

func calculateStockScore(rec domain.Recommendation) float64 {
	targetRange := rec.TargetTo - rec.TargetFrom
	if targetRange < 0 {
		targetRange = 0
	}

	averageTarget := (rec.TargetFrom + rec.TargetTo) / 2

	upgradeScore := 0.0
	if rec.Action == "upgraded by" {
		upgradeScore = 10.0
	}

	score := targetRange*10 + averageTarget + upgradeScore
	return score
}
