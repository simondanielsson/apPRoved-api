package main

import (
	"log"

	"github.com/simondanielsson/apPRoved/cmd/api"
	"github.com/simondanielsson/apPRoved/cmd/config"
	"github.com/simondanielsson/apPRoved/cmd/internal/db"
	"github.com/simondanielsson/apPRoved/pkg/utils"
	"github.com/simondanielsson/apPRoved/pkg/utils/mq"
)

// @title apPRoved API
// @version 1.0
// @description API for apPRoved
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config %v\n", err)
	}

	utils.SetJWTKey(config.JWT.Secret)

	db, err := db.NewDB(config.Database)
	if err != nil {
		log.Fatalf("could not connect to database %v\n", err)
	}

	if err := mq.InitRabbitMQ(config.MQ); err != nil {
		log.Fatalf("could not connect to RabbitMQ")
	}
	defer mq.CloseRabbitMQ()

	server := api.NewAPIServer(config.Server, db)
	server.Run()
}
