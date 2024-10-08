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
	"github.com/simondanielsson/apPRoved/pkg/utils/mq"
)

type ReviewsController struct {
	reviewsService *services.ReviewsService
}

// NewReviewsController creates a new reviews controller
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

	githubClient, ok := c.Locals("githubClient").(utils.GithubClient)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Github client not available")
	}

	repo, err := rc.reviewsService.RegisterRepository(ctx, tx, githubClient, userID, req.Name, req.Owner, req.URL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Could not create repository", "error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Repository created",
		"data":    repo,
	})
}

// @Summary Get repository
// @Description Get a repository
// @Tags reviews
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param        repositoryID  path  string  true  "Repository ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/repositories/{repositoryID} [get]
func (rc *ReviewsController) GetRepository(c *fiber.Ctx) error {
	repoID, err := utils.ReadUintPathParam(c, "repositoryID")
	if err != nil {
		return err
	}

	tx := db.GetDBTransaction(c)

	repo, err := rc.reviewsService.GetRepository(tx, repoID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Repository not found", "error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Repository fetched successfully",
		"data":    repo,
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

// generate swagger docs
// @Summary Update pull request
// @Description Update a pull request
// @Tags reviews
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param        repositoryID  path  string  true  "Repository ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/repositories/{repositoryID}/pull-requests [put]
func (rc *ReviewsController) RefreshPullRequests(c *fiber.Ctx) error {
	repoID, err := utils.ReadUintPathParam(c, "repositoryID")
	if err != nil {
		return err
	}
	userID := middlewares.GetUserID(c)

	githubClient, ok := c.Locals("githubClient").(utils.GithubClient)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Github client not available")
	}

	tx := db.GetDBTransaction(c)
	ctx := context.Background()
	if err := rc.reviewsService.RefreshPullRequests(ctx, tx, githubClient, userID, repoID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not update pull requests",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully refreshed pull requests",
	})
}

// @Summary Get pull request
// @Description Get a pull request
// @Tags reviews
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param        repositoryID  path  string  true  "Repository ID"
// @Param        prID          path  string  true  "Pull request ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/repositories/{repositoryID}/pull-requests/{prID} [get]
func (rc *ReviewsController) GetPullRequest(c *fiber.Ctx) error {
	repoID, err := utils.ReadUintPathParam(c, "repositoryID")
	if err != nil {
		return err
	}

	prID, err := utils.ReadUintPathParam(c, "prID")
	if err != nil {
		return err
	}

	tx := db.GetDBTransaction(c)
	pr, err := rc.reviewsService.GetPullRequest(tx, repoID, prID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not fetch pull request",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully fetched pull request",
		"data":    pr,
	})
} // generate swagger docs
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
	repoID, err := utils.ReadUintPathParam(c, "repositoryID")
	if err != nil {
		return err
	}
	prID, err := utils.ReadUintPathParam(c, "prID")
	if err != nil {
		return err
	}
	reviewID, err := utils.ReadUintPathParam(c, "reviewID")
	if err != nil {
		return err
	}

	tx := db.GetDBTransaction(c)
	review, err := rc.reviewsService.GetReview(tx, repoID, prID, reviewID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not fetch review",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully fetched review",
		"data":    review,
	})
}

// generate swagger docs
// @Summary Get file reviews
// @Description Get all file reviews for a review
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
// @Router       /api/v1/repositories/{repositoryID}/pull-requests/{prID}/reviews/{reviewID}/files [get]
func (rc *ReviewsController) GetFileReviews(c *fiber.Ctx) error {
	reviewID, err := utils.ReadUintPathParam(c, "reviewID")
	if err != nil {
		return err
	}

	tx := db.GetDBTransaction(c)
	reviewResponse, err := rc.reviewsService.GetFileReviews(tx, reviewID)
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
	messageQueue, ok := c.Locals("messageQueue").(mq.MessageQueue)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Message queue not available")
	}
	githubClient, ok := c.Locals("githubClient").(utils.GithubClient)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Github client not available")
	}

	ctx := context.Background()
	review, err := rc.reviewsService.CreateReview(tx, ctx, messageQueue, githubClient, repoID, prID, req.Name, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not create review",
			"error":   err.Error(),
		})
	}

	c.Set("Location", fmt.Sprintf("/api/v1/repositories/%d/pull-requests/%d/reviews/%d/progress", repoID, prID, review.ID))

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Review initiated.",
		"data":    review,
	})
}

// Generate swagger docs
// @Summary Delete review
// @Description Delete a review
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
// @Router       /api/v1/repositories/{repositoryID}/pull-requests/{prID}/reviews/{reviewID} [delete]
func (rc *ReviewsController) DeleteReview(c *fiber.Ctx) error {
	repoID, err := utils.ReadUintPathParam(c, "repositoryID")
	if err != nil {
		return err
	}
	prID, err := utils.ReadUintPathParam(c, "prID")
	if err != nil {
		return err
	}
	reviewID, err := utils.ReadUintPathParam(c, "reviewID")
	if err != nil {
		return err
	}

	tx := db.GetDBTransaction(c)
	if err := rc.reviewsService.DeleteReview(tx, repoID, prID, reviewID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not delete review",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Review deleted successfully",
	})
}

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

// write swagger docs
// @Summary Get review progress
// @Description Get review progress
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
// @Router       /api/v1/repositories/{repositoryID}/pull-requests/{prID}/reviews/{reviewID}/progress [get]
func (rc *ReviewsController) GetReviewProgress(c *fiber.Ctx) error {
	repoID, err := utils.ReadUintPathParam(c, "repositoryID")
	if err != nil {
		return err
	}
	prID, err := utils.ReadUintPathParam(c, "prID")
	if err != nil {
		return err
	}
	reviewID, err := utils.ReadUintPathParam(c, "reviewID")
	if err != nil {
		return err
	}

	tx := db.GetDBTransaction(c)
	reviewStatus, err := rc.reviewsService.GetReviewStatus(tx, repoID, prID, reviewID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not get review progress",
			"error":   err.Error(),
		})
	}

	c.Set("Location", fmt.Sprintf("/api/v1/repositories/%d/pull-requests/%d/reviews/%d/progress", repoID, prID, reviewID))

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Fetched progress successfully.",
		"data": fiber.Map{
			"status":   reviewStatus.Status,
			"progress": reviewStatus.Progress,
		},
	})
}

// generate swagger docs
// @Summary Update review progress
// @Description Update review progress
// @Tags reviews
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param        reviewStatusID  path  string  true  "Review Status ID"
// @Param        updateReviewRequest  body  requests.UpdateReviewRequest  true  "Update review progress request"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/review-status/{reviewStatusID} [put]
func (rc *ReviewsController) UpdateReviewProgress(c *fiber.Ctx) error {
	reviewStatusID, err := utils.ReadUintPathParam(c, "reviewStatusID")
	if err != nil {
		return err
	}
	var request requests.UpdateReviewRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Could not parse request body"})
	}

	tx := db.GetDBTransaction(c)
	if err := rc.reviewsService.UpdateReviewStatus(tx, reviewStatusID, request.Status, request.Progress); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not update review progress",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Updated progress successfully.",
	})
}
