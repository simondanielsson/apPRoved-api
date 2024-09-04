package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/simondanielsson/apPRoved/cmd/internal/controllers"
	"github.com/simondanielsson/apPRoved/cmd/internal/middlewares"
)

// @Summary Health check
// @Description Check if the service is up and running
// @Tags health
// @Accept json
// @Produce json
// @Router       /api/v1/health [get]
func Health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}

func RegisterRoutes(apiV1 fiber.Router, ctrls *controllers.Controllers, opt_middlewares middlewares.OptionalMiddlewares) {
	apiV1.Get("/health", Health)

	RegisterAuthRoutes(apiV1, ctrls.AuthController, opt_middlewares)
	RegisterReviewsRoutes(apiV1, ctrls.ReviewsController, opt_middlewares)
	RegisterUserRoutes(apiV1, ctrls.UserController, opt_middlewares)
}
