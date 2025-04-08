package services

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"stock-information/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

func InitDatabase() *gorm.DB {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	port := os.Getenv("DB_PORT")
	ssl := os.Getenv("DB_SSL")
	dbName := os.Getenv("DB_NAME")

	tempDsn := fmt.Sprintf("host=%s user=%s dbname=defaultdb port=%s sslmode=%s",
		host, user, port, ssl,
	)

	tempDB, err := sql.Open("postgres", tempDsn)
	if err != nil {
		log.Fatalf("❌ Error conectando a defaultdb: %v", err)
	}
	defer tempDB.Close()

	_, err = tempDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", dbName))
	if err != nil {
		log.Fatalf("❌ Error al crear la base de datos %s: %v", dbName, err)
	}
	log.Printf("✅ Base de datos '%s' verificada o creada correctamente", dbName)

	var gormDB *gorm.DB
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=%s",
		host, user, dbName, port, ssl,
	)

	for i := 0; i < 10; i++ {
		gormDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("⏳ Intentando conexión con la base de datos '%s'... (%d/10)", dbName, i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("❌ No se pudo conectar a CockroachDB: %v", err)
	}

	err = gormDB.AutoMigrate(&domain.Recommendation{})
	if err != nil {
		log.Fatalf("❌ Error en AutoMigrate: %v", err)
	}

	log.Println("✅ Migración completada y conexión establecida con éxito 🎉")
	return gormDB
}
