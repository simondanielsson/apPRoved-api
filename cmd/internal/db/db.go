package db

import (
	"fmt"
	"log"

	"github.com/simondanielsson/apPRoved/cmd/config"
	"github.com/simondanielsson/apPRoved/cmd/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	log.Printf("connecting to database %s\n", dsn)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DriverName: "pgx",
		DSN:        dsn,
	}), &gorm.Config{})
	if err != nil {
		log.Fatalf("could not connect to database %v\n", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Repository{}); err != nil {
		log.Fatalf("could not migrate database %v\n", err)
	}

	log.Print("connected to database\n")
	return db, nil
}
