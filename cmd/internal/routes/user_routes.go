package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/simondanielsson/apPRoved/cmd/internal/controllers"
	"github.com/simondanielsson/apPRoved/cmd/internal/middlewares"
)

func RegisterUserRoutes(apiV1 fiber.Router, userController *controllers.UserController, opt_middlewares middlewares.OptionalMiddlewares) {
	router := apiV1.Group("/users", opt_middlewares.Auth, opt_middlewares.Transaction)

	router.Get("", userController.GetUsers)
	router.Get("/:id", userController.GetUser)
	router.Delete("/:id", userController.DeleteUser)
}
