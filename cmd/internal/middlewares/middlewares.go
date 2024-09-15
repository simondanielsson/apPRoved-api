package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/gorm"
)

type OptionalMiddlewares struct {
	Auth        func(*fiber.Ctx) error
	Transaction func(*fiber.Ctx) error
}

func SetupMiddlewares(app *fiber.App) {
	app.Use(cors.New(
		cors.Config{
			AllowOrigins:     "http://localhost:5173, http://host.docker.internal:5173, http://approved-frontend:5173",
			AllowCredentials: true,
			AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
			AllowHeaders:     "Authorization,Content-Type",
		},
	))

	app.Use(logger.New(logger.Config{
		Format:     "${time} ${status} - ${method} ${path}\n",
		TimeFormat: "2 Jan 2006 15:04:05",
		TimeZone:   "local",
	}))
}

func GetOptionalMiddlewares(db *gorm.DB) OptionalMiddlewares {
	return OptionalMiddlewares{
		Transaction: GetTransactionMiddleware(db),
		Auth:        AuthMiddleware,
	}
}
