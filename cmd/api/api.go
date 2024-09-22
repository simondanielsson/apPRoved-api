package api

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/simondanielsson/apPRoved/cmd/bootstrap"
	"github.com/simondanielsson/apPRoved/cmd/config"
	"github.com/simondanielsson/apPRoved/cmd/internal/middlewares"
	"github.com/simondanielsson/apPRoved/cmd/internal/routes"
	"github.com/simondanielsson/apPRoved/pkg/utils"
	"github.com/simondanielsson/apPRoved/pkg/utils/mq"
	"gorm.io/gorm"
)

type APIServer struct {
	config       *config.ServerConfig
	db           *gorm.DB
	app          *fiber.App
	queue        mq.MessageQueue
	githubClient *utils.GithubClient
}

func (s *APIServer) Run() {
	if err := s.app.Listen(":" + s.config.BindAddr); err != nil {
		log.Fatalf("could not start server %v\n", err)
	}
	log.Printf("API server listening on port %s", s.config.BindAddr)
}

func (s *APIServer) setupRoutes() {
	apiV1 := s.app.Group("/api/v1")

	repos := bootstrap.InitRepositories()
	services := bootstrap.InitServices(repos)
	controllers := bootstrap.InitControllers(services)

	opt_middlewares := middlewares.GetOptionalMiddlewares(s.db)
	routes.RegisterRoutes(apiV1, controllers, opt_middlewares)
}

func (s *APIServer) Shutdown() error {
	fmt.Println("Shutting down server...")
	if err := s.app.Shutdown(); err != nil {
		return err
	}

	return nil
}

func NewAPIServer(cfg *config.ServerConfig, db *gorm.DB, queue mq.MessageQueue, githubClient *utils.GithubClient) *APIServer {
	server := &APIServer{
		config:       cfg,
		db:           db,
		app:          fiber.New(),
		queue:        queue,
		githubClient: githubClient,
	}

	utils.ConfigureSwagger(server.app)
	middlewares.SetupMiddlewares(server.app, queue, githubClient)
	server.setupRoutes()

	return server
}
