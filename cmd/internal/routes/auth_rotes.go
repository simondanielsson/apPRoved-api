package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/simondanielsson/apPRoved/cmd/internal/controllers"
)

func RegisterAuthRoutes(apiV1 fiber.Router, authController *controllers.AuthController) {
	apiV1.Post("/login", authController.Login)
}
