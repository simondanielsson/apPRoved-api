package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	_ "github.com/simondanielsson/apPRoved/docs"
)

func ConfigureSwagger(app *fiber.App) {
	app.Get("/swagger/*", swagger.HandlerDefault)
}
