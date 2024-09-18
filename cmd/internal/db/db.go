package db

import (
	"fmt"
	"log"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/gofiber/fiber/v2"
	"github.com/simondanielsson/apPRoved/cmd/config"
	"github.com/simondanielsson/apPRoved/cmd/internal/middlewares"
	"github.com/simondanielsson/apPRoved/cmd/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	var dsn string
	switch cfg.DriverName {
	case "pgx":
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	case "cloudsqlpostgres":
		dsn = fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable", cfg.Host, cfg.User, cfg.DBName, cfg.Password)
	default:
		log.Fatalf("unsupported driver name: %s", cfg.DriverName)
	}

	log.Printf("connecting to database %s:%s\n", cfg.Host, cfg.DBName)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DriverName: cfg.DriverName,
		DSN:        dsn,
	}), &gorm.Config{})
	if err != nil {
		log.Fatalf("could not connect to database %v\n", err)
	}

	underlyingDB, err := db.DB()
	if err != nil {
		log.Fatalf("could not get underlying db connection: %v\n", err)
	}
	if err := underlyingDB.Ping(); err != nil {
		log.Fatalf("could not ping database: %v\n", err)
	}

	for _, model := range models.Models {
		if err := db.AutoMigrate(model); err != nil {
			log.Fatalf("failed to migrate model: %v", err)
		}
	}
	log.Print("connected to database\n")
	return db, nil
}

func GetDBTransaction(c *fiber.Ctx) *gorm.DB {
	tx, ok := c.Locals(string(middlewares.TxnKey)).(*gorm.DB)
	if !ok {
		panic("no transaction found in context, did you forget to wrap the route in a transaction middleware?")
	}
	return tx
}
