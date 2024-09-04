package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/simondanielsson/apPRoved/pkg/utils"
)

// AuthMiddleware is a middleware that checks if the request has a valid JWT token, and attaches the user ID to the context.
func AuthMiddleware(c *fiber.Ctx) error {
	// Authorization: Bearer <token>
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "missing Authorization header"})
	}

	token := authHeader[len("Bearer "):]
	claims, err := utils.ParseJWTToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid token"})
	}

	c.Locals("userID", claims.UserID)
	return c.Next()
}
