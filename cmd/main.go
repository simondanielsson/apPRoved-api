package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

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

	messageQueue, err := mq.NewMessageQueue(config)
	if err != nil {
		log.Fatalf("could not connect to message queue: %v", err)
	}
	defer messageQueue.Close()

	server := api.NewAPIServer(config.Server, db, messageQueue)

	gracefulShutdown(server, &messageQueue)
	server.Run()
}

func gracefulShutdown(server *api.APIServer, queue *mq.MessageQueue) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-quit
		log.Println("Shutting down server...")

		(*queue).Close()

		if err := server.Shutdown(); err != nil {
			log.Fatalf("Failed to shutdown server: %v", err)
		}

		log.Println("Server stopped gracefully")
	}()
}
