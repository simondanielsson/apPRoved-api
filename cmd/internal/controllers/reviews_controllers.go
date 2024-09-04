package controllers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/simondanielsson/apPRoved/cmd/internal/dto/requests"
	"github.com/simondanielsson/apPRoved/cmd/internal/services"
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
// @Router       /api/v1/reviews/repositories [get]
func (rc *ReviewsController) GetRepositories(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	repos, err := rc.reviewsService.GetRepositories(userID)
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
// @Router       /api/v1/reviews/repositories [post]
func (rc *ReviewsController) RegisterRepository(c *fiber.Ctx) error {
	ctx := context.Background()
	userID := c.Locals("userID").(uint)

	var req requests.CreateRepositoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Could not parse request body"})
	}

	repoID, err := rc.reviewsService.RegisterRepository(ctx, userID, req.Name, req.Owner, req.URL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Could not create repository", "error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Repository created",
		"id":      repoID,
	})
}
