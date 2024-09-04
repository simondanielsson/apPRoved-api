package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/simondanielsson/apPRoved/cmd/internal/dto/requests"
	"github.com/simondanielsson/apPRoved/cmd/internal/dto/responses"
	"github.com/simondanielsson/apPRoved/cmd/internal/services"
)

type AuthController struct {
	authService *services.AuthService
}

// NewAuthController creates a new auth controller
func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// generate swagger docs for this; exoects username and paossword
// @Summary Login
// @Description Login to the application
// @Tags auth
// @Accept json
// @Produce json
// @Param        loginRequest  body      requests.LoginRequest  true  "User login request"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {string}  string "Unauthorized"
// @Router       /api/v1/login [post]
func (ac *AuthController) Login(ctx *fiber.Ctx) error {
	var loginRequest requests.LoginRequest

	if err := ctx.BodyParser(&loginRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Could not parse request body", "error": err.Error()})
	}

	token, err := ac.authService.AuthenticateUser(loginRequest.Username, loginRequest.Password)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid credentials", "error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.LoginResponse{Token: token})
}
