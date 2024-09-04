package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/simondanielsson/apPRoved/cmd/internal/controllers"
	"github.com/simondanielsson/apPRoved/cmd/internal/middlewares"
)

func RegisterUserRoutes(apiV1 fiber.Router, userController *controllers.UserController) {
	router := apiV1.Group("/users", middlewares.AuthMiddleware)

	router.Get("", userController.GetUsers)
	router.Get("/:id", userController.GetUser)
	router.Post("", userController.CreateUser)
	router.Delete("/:id", userController.DeleteUser)
}
