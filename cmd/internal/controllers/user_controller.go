package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/simondanielsson/apPRoved/cmd/internal/dto/requests"
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
	users, err := uc.userService.GetUsers()
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

	user, err := uc.userService.GetUser(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Could not fetch user", "error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"user": *user})
}

// CreateUser godoc
// @Summary      Create a new user
// @Description  Create a new user with a name, email, and password
// @Tags         users
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        user  body      requests.CreateUserRequest  true  "User creation request"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/users [post]
func (uc *UserController) CreateUser(c *fiber.Ctx) error {
	var req requests.CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Could not parse request body"})
	}

	userID, err := uc.userService.CreateUser(req.Username, req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Could not create user", "error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created",
		"id":      userID,
	})
}

// DeleteUser deletes a user
func (uc *UserController) DeleteUser(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"message": "Not implemented"})
}
