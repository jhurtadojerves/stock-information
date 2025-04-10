package main

import (
	"fmt"
	"log"
	"stock-information/internal/adapters/repository"
	"stock-information/internal/controllers"
	"stock-information/internal/services"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	gormDB := services.InitDatabase()
	repo := &repository.GormCockroachRepo{DB: gormDB}
	stockService := &services.StockService{Repo: repo}
	stockController := &controllers.StockController{StockService: stockService}
	scheduler := services.Scheduler{Repo: repo}
	scheduler.Start()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4321"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	v1 := r.Group("/api/v1")
	stocks := v1.Group("/stocks")
	{
		stocks.GET("", stockController.GetAllRecommendations)
		stocks.GET("/recommendation", stockController.GetBestRecommendation)
		stocks.GET("/:ticker", stockController.GetRecommendationByTicker)
	}
	port := "8080"
	log.Printf("🌐 Servidor HTTP escuchando en :%s", port)
	err := r.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("❌ Error iniciando el servidor HTTP: %v", err)
	}
}
