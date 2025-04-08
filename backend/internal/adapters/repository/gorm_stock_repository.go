package repository

import (
	"stock-information/internal/domain"
	"stock-information/internal/ports"
	"strings"

	"gorm.io/gorm"
)

type GormCockroachRepo struct {
	DB *gorm.DB
}

func (r *GormCockroachRepo) SaveRecommendations(recs []domain.Recommendation) error {
	for _, rec := range recs {
		err := r.DB.Save(&rec).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *GormCockroachRepo) GetRecommendations(page, limit int, sort string) ([]domain.Recommendation, int, error) {
	var recs []domain.Recommendation
	var totalCount int64

	err := r.DB.Model(&domain.Recommendation{}).Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	var query *gorm.DB

	if sort == "desc" {
		query = r.DB.Offset((page - 1) * limit).Limit(limit).Order("ticker desc")
	} else {
		query = r.DB.Offset((page - 1) * limit).Limit(limit).Order("ticker asc")
	}

	err = query.Find(&recs).Error
	return recs, int(totalCount), err
}

func (r *GormCockroachRepo) GetRecommendationsByCompany(company string, page, limit int, sort string) ([]domain.Recommendation, int, error) {
	var recs []domain.Recommendation
	var totalCount int64

	company = strings.ToLower(company)

	err := r.DB.Model(&domain.Recommendation{}).Where("LOWER(company) LIKE ?", "%"+company+"%").Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	var query *gorm.DB
	if sort == "desc" {
		query = r.DB.Where("LOWER(company) LIKE ?", "%"+company+"%").Offset((page - 1) * limit).Limit(limit).Order("ticker desc")
	} else {
		query = r.DB.Where("LOWER(company) LIKE ?", "%"+company+"%").Offset((page - 1) * limit).Limit(limit).Order("ticker asc")
	}

	err = query.Find(&recs).Error
	return recs, int(totalCount), err
}

func (r *GormCockroachRepo) GetRecommendationByTicker(ticker string) (domain.Recommendation, error) {
	var rec domain.Recommendation
	err := r.DB.Where("ticker = ?", ticker).First(&rec).Error
	if err != nil {
		return domain.Recommendation{}, err
	}
	return rec, nil
}

func (r *GormCockroachRepo) GetAllRecommendations() ([]domain.Recommendation, error) {
	var recs []domain.Recommendation
	err := r.DB.Find(&recs).Error
	return recs, err
}

var _ ports.StockRepository = (*GormCockroachRepo)(nil)
