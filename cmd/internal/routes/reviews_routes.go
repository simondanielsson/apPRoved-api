package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/simondanielsson/apPRoved/cmd/internal/controllers"
	"github.com/simondanielsson/apPRoved/cmd/internal/middlewares"
)

func RegisterReviewsRoutes(apiV1 fiber.Router, reviewsController *controllers.ReviewsController, opt_middlewares middlewares.OptionalMiddlewares) {
	router := apiV1.Group("/reviews", opt_middlewares.Auth, opt_middlewares.Transaction)

	router.Get("/repositories", reviewsController.GetRepositories)
	router.Post("/repositories", reviewsController.RegisterRepository)

	// router.Get("/repositories/:repositoryID/pull-requests", reviewsController.GetPullRequests)

	// router.Get("/repositories/:repositoryID/pull-requests/:prID/reviews", reviewsController.GetReviews)
	// router.Get("/repositories/:repositoryID/pull-requests/:prID/reviews/:reviewID", reviewsController.GetReview)
	// router.Post("/repositories/:repositoryID/pull-requests/:prID/reviews", reviewsController.CreateReview)
}
