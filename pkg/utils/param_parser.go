package utils

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// ReadUintParam reads an unsigned integer parameter from the request.
func ReadUintPathParam(c *fiber.Ctx, name string) (uint, error) {
	value := c.Params(name)

	if value == "" {
		return 0, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": fmt.Sprintf("%s is required", name),
		})
	}

	valueUint, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": fmt.Sprintf("Invalid %s, got %s. Should be an unsigned integer.", name, value),
		})
	}
	return uint(valueUint), nil
}
