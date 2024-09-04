package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/simondanielsson/apPRoved/cmd/internal/controllers"
	"github.com/simondanielsson/apPRoved/cmd/internal/middlewares"
)

func RegisterReviewsRoutes(apiV1 fiber.Router, reviewsController *controllers.ReviewsController) {
	router := apiV1.Group("/reviews", middlewares.AuthMiddleware)

	router.Get("/repositories", reviewsController.GetRepositories)
	router.Post("/repositories", reviewsController.CreateRepository)
}
