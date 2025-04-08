package services

import (
	"log"
	"os"

	"stock-information/internal/ports"
	"stock-information/external"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	Repo ports.StockRepository
}

func (s *Scheduler) Start() {
	cronExpr := os.Getenv("SYNC_CRON_EXPRESSION")
	if cronExpr == "" {
		cronExpr = "0 * * * *" // valor por defecto: cada hora
		log.Printf("ℹ️  SYNC_CRON_EXPRESSION no definida, usando valor por defecto: %s", cronExpr)
	} else {
		log.Printf("🕒 Usando expresión CRON desde env: %s", cronExpr)
	}

	c := cron.New()

	_, err := c.AddFunc(cronExpr, func() {
		log.Println("⏰ Ejecutando sincronización programada...")

		recs, err := external.FetchAllRecommendations()
		if err != nil {
			log.Printf("❌ Error al obtener datos del API externo: %v", err)
			return
		}

		service := StockService{Repo: s.Repo}
		if err := service.SyncRecommendations(recs); err != nil {
			log.Printf("❌ Error al guardar recomendaciones: %v", err)
			return
		}

		log.Printf("✅ Sincronización completada (%d registros)", len(recs))
	})

	if err != nil {
		log.Fatalf("❌ Error al registrar tarea cron: %v", err)
	}

	c.Start()
	log.Println("📆 Scheduler iniciado.")
}
