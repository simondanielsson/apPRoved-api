package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/simondanielsson/apPRoved/cmd/internal/db"
	"github.com/simondanielsson/apPRoved/cmd/internal/dto/responses"
	"github.com/simondanielsson/apPRoved/cmd/internal/services"
)

type AuthController struct {
	authService *services.AuthService
	userService *services.UserService
}

// NewAuthController creates a new auth controller
func NewAuthController(authService *services.AuthService, userService *services.UserService) *AuthController {
	return &AuthController{
		authService: authService,
		userService: userService,
	}
}

// generate swagger docs for this; exoects username and paossword
// @Summary Login
// @Description Login to the application
// @Tags auth
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {string} string "Unauthorized"
// @Router       /api/v1/login [post]
func (ac *AuthController) Login(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "" || password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Username and password are required"})
	}

	tx := db.GetDBTransaction(c)
	token, err := ac.authService.AuthenticateUser(tx, username, password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid credentials", "error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(responses.LoginResponse{Token: token})
}

// CreateUser godoc
// @Summary      Register
// @Description  Register a new user with a name, email, and password
// @Tags         auth
// @Security BearerAuth
// @Accept application/x-www-form-urlencoded
// @Produce      json
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Param email formData string true "Email"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/register [post]
func (ac *AuthController) Register(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	email := c.FormValue("email")

	if username == "" || password == "" || email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Username, email and password are required"})
	}

	tx := db.GetDBTransaction(c)
	userID, err := ac.userService.CreateUser(tx, username, email, password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Could not create user", "error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created",
		"id":      userID,
	})
}
