package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/simondanielsson/apPRoved/cmd/internal/controllers"
	"github.com/simondanielsson/apPRoved/cmd/internal/middlewares"
)

func RegisterReviewsRoutes(apiV1 fiber.Router, reviewsController *controllers.ReviewsController, opt_middlewares middlewares.OptionalMiddlewares) {
	router := apiV1.Group("/repositories", opt_middlewares.Auth, opt_middlewares.Transaction)

	router.Get("", reviewsController.GetRepositories)
	router.Post("", reviewsController.RegisterRepository)

	router.Get("/:repositoryID/pull-requests", reviewsController.GetPullRequests)
	router.Put("/:repositoryID/pull-requests/:prID", reviewsController.UpdatePullRequest)

	router.Get("/:repositoryID/pull-requests/:prID/reviews", reviewsController.GetReviews)
	router.Get("/:repositoryID/pull-requests/:prID/reviews/:reviewID", reviewsController.GetReview)
	router.Post("/:repositoryID/pull-requests/:prID/reviews", reviewsController.CreateReview)
	router.Get("/:repositoryID/pull-requests/:prID/reviews/:reviewID/progress", reviewsController.GetReviewProgress)
	router.Delete("/:repositoryID/pull-requests/:prID/reviews/:reviewID", reviewsController.DeleteReview)

	apiV1.Post("/reviews/complete", opt_middlewares.Transaction, reviewsController.CompleteReview)
}
