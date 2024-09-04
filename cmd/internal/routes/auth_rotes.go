package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/simondanielsson/apPRoved/cmd/internal/controllers"
	"github.com/simondanielsson/apPRoved/cmd/internal/middlewares"
)

func RegisterAuthRoutes(apiV1 fiber.Router, authController *controllers.AuthController, opt_middlewares middlewares.OptionalMiddlewares) {
	apiV1.Post("/login", opt_middlewares.Transaction, authController.Login)
}
