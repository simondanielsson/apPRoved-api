package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/simondanielsson/apPRoved/cmd/internal/db"
	"github.com/simondanielsson/apPRoved/cmd/internal/services"
)

// UserController handles user-related endpoints
type UserController struct {
	userService *services.UserService
}

// NewUserController creates a new instance of UserController
func NewUserController(userService *services.UserService) *UserController {
	return &UserController{userService: userService}
}

// GetUsers returns a list of users
// @Summary      Get a list of users
// @Description  Get a list of all users
// @Tags         users
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/users [get]
func (uc *UserController) GetUsers(c *fiber.Ctx) error {
	tx := db.GetDBTransaction(c)
	users, err := uc.userService.GetUsers(tx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Could not fetch users", "error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"users": users})
}

// create swagger documentation for GetUser
// GetUser returns a user by ID
// @Summary      Get a user by ID
// @Description  Get a user by ID
// @Tags         users
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "User ID"
// @Router       /api/v1/users/{id} [get]
func (uc *UserController) GetUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID. Should be a positive integer"})
	}
	if id < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID. Should be a positive integer"})
	}
	userID := uint(id)

	tx := db.GetDBTransaction(c)
	user, err := uc.userService.GetUser(tx, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Could not fetch user", "error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"user": *user})
}

// DeleteUser deletes a user
func (uc *UserController) DeleteUser(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"message": "Not implemented"})
}
