package controllers

import (
	"stock-information/internal/domain"
	"stock-information/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StockController struct {
	StockService *services.StockService
}

func (controller *StockController) GetAllRecommendations(c *gin.Context) {
	company := c.DefaultQuery("company", "")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sort := c.DefaultQuery("sort", "asc")

	var recommendations []domain.Recommendation
	var totalCount int
	var err error

	if company == "" {
		recommendations, totalCount, err = controller.StockService.GetRecommendations(page, limit, sort)
	} else {
		recommendations, totalCount, err = controller.StockService.GetRecommendationsByCompany(company, page, limit, sort)
	}

	if err != nil {
		c.JSON(500, gin.H{"message": "Error obteniendo los datos"})
		return
	}

	nextPage := page
	if totalCount > page*limit {
		nextPage = page + 1
	}

	c.JSON(200, gin.H{
		"data":        recommendations,
		"nextPage":    nextPage,
		"total":       totalCount,
		"totalInPage": len(recommendations),
		"page":        page,
		"limit":       limit,
		"sort":        sort,
	})
}

func (controller *StockController) GetBestRecommendation(c *gin.Context) {
	recommendation, err := controller.StockService.GetBestRecommendation()
	if err != nil {
		c.JSON(500, gin.H{"message": "Error recomendando la mejor acción"})
		return
	}
	c.JSON(200, recommendation)
}

func (controller *StockController) GetRecommendationByTicker(c *gin.Context) {
	ticker := c.Param("ticker")

	recommendation, err := controller.StockService.GetRecommendationByTicker(ticker)
	if err != nil {
		c.JSON(404, gin.H{"message": "Recomendación no encontrada para el ticker " + ticker})
		return
	}

	c.JSON(200, recommendation)
}
