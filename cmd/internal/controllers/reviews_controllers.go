package controllers

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/simondanielsson/apPRoved/cmd/internal/db"
	"github.com/simondanielsson/apPRoved/cmd/internal/dto/requests"
	"github.com/simondanielsson/apPRoved/cmd/internal/middlewares"
	"github.com/simondanielsson/apPRoved/cmd/internal/services"
	"github.com/simondanielsson/apPRoved/pkg/utils"
)

type ReviewsController struct {
	reviewsService *services.ReviewsService
}

func NewReviewsController(reviewsService *services.ReviewsService) *ReviewsController {
	return &ReviewsController{reviewsService: reviewsService}
}

// generate swagger docs
// @Summary Get repositories
// @Description Get all repositories for a user
// @Tags reviews
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/repositories [get]
func (rc *ReviewsController) GetRepositories(c *fiber.Ctx) error {
	tx := db.GetDBTransaction(c)
	userID := middlewares.GetUserID(c)

	repos, err := rc.reviewsService.GetRepositories(tx, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No repositories found", "error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"repositories": repos,
	})
}

// @Summary Create repository
// @Description Create a new repository
// @Tags reviews
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param        createRepositoryRequest  body      requests.CreateRepositoryRequest  true  "Create repository request"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/repositories [post]
func (rc *ReviewsController) RegisterRepository(c *fiber.Ctx) error {
	ctx := context.Background()
	tx := db.GetDBTransaction(c)
	userID := middlewares.GetUserID(c)

	var req requests.CreateRepositoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Could not parse request body"})
	}

	repoID, err := rc.reviewsService.RegisterRepository(ctx, tx, userID, req.Name, req.Owner, req.URL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Could not create repository", "error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Repository created",
		"id":      repoID,
	})
}

// generate swagger docs
// @Summary Get pull requests
// @Description Get all pull requests for repository
// @Tags reviews
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param        repositoryID  path  string  true  "Repository ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/repositories/{repositoryID}/pull-requests [get]
func (rc *ReviewsController) GetPullRequests(c *fiber.Ctx) error {
	repoID, err := utils.ReadUintPathParam(c, "repositoryID")
	if err != nil {
		return err
	}
	userID := middlewares.GetUserID(c)

	tx := db.GetDBTransaction(c)
	prs, err := rc.reviewsService.GetPullRequests(tx, userID, repoID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not fetch pull requests",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully fetched pull requests",
		"data":    prs,
	})
}

func (rc *ReviewsController) UpdatePullRequest(c *fiber.Ctx) error {
	return c.SendString("UpdatePullRequest")
}

// generate swagger docs
// @Summary Get reviews
// @Description Get all reviews for a pull request
// @Tags reviews
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param        repositoryID  path  string  true  "Repository ID"
// @Param        prID          path  string  true  "Pull request ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/repositories/{repositoryID}/pull-requests/{prID}/reviews [get]
func (rc *ReviewsController) GetReviews(c *fiber.Ctx) error {
	repoID, err := utils.ReadUintPathParam(c, "repositoryID")
	if err != nil {
		return err
	}
	prID, err := utils.ReadUintPathParam(c, "prID")
	if err != nil {
		return err
	}

	tx := db.GetDBTransaction(c)
	reviewsResponse, err := rc.reviewsService.GetReviews(tx, repoID, prID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not fetch reviews",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully fetched reviews",
		"data":    reviewsResponse,
	})
}

// generate swagger docs
// @Summary Get review
// @Description Get a review
// @Tags reviews
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param        repositoryID  path  string  true  "Repository ID"
// @Param        prID          path  string  true  "Pull request ID"
// @Param        reviewID  path  string  true  "Review ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/repositories/{repositoryID}/pull-requests/{prID}/reviews/{reviewID} [get]
func (rc *ReviewsController) GetReview(c *fiber.Ctx) error {
	reviewID, err := utils.ReadUintPathParam(c, "reviewID")
	if err != nil {
		return err
	}

	tx := db.GetDBTransaction(c)
	reviewResponse, err := rc.reviewsService.GetReview(tx, reviewID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not fetch review",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully fetched review",
		"data":    reviewResponse,
	})
}

// generate swagger docs
// @Summary Create review
// @Description Create a review
// @Tags reviews
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param        repositoryID  path  string  true  "Repository ID"
// @Param        prID          path  string  true  "Pull request ID"
// @Param        createReviewRequest  body  requests.CreateReviewRequest  true  "Create review request"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/repositories/{repositoryID}/pull-requests/{prID}/reviews [post]
func (rc *ReviewsController) CreateReview(c *fiber.Ctx) error {
	repoID, err := utils.ReadUintPathParam(c, "repositoryID")
	if err != nil {
		return err
	}
	prID, err := utils.ReadUintPathParam(c, "prID")
	if err != nil {
		return err
	}
	var req requests.CreateReviewRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Could not parse request body"})
	}

	tx := db.GetDBTransaction(c)
	userID := middlewares.GetUserID(c)
	ctx := context.Background()
	reviewID, err := rc.reviewsService.CreateReview(tx, ctx, repoID, prID, req.Name, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not create review",
			"error":   err.Error(),
		})
	}

	c.Set("Location", fmt.Sprintf("/api/v1/repositories/%d/pull-requests/%d/reviews/%d/status", repoID, prID, reviewID))

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Review initiated.",
		"id":      reviewID,
	})
}

// generate swagger docs
// @Summary Complete review
// @Description Complete a review
// @Tags reviews
// @Accept json
// @Produce json
// @Param        completeReviewRequest  body  requests.CompleteReviewRequest  true  "Update review status request"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/reviews/complete [post]
func (rc *ReviewsController) CompleteReview(c *fiber.Ctx) error {
	var req requests.CompleteReviewRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Could not parse request body"})
	}

	tx := db.GetDBTransaction(c)
	if err := rc.reviewsService.CompleteReview(tx, &req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not complete review",
			"error":   err.Error,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Review completed.",
	})
}
