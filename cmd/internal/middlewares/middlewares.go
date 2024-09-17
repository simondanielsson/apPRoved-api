package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/simondanielsson/apPRoved/pkg/utils/mq"
	"gorm.io/gorm"
)

type OptionalMiddlewares struct {
	Auth        func(*fiber.Ctx) error
	Transaction func(*fiber.Ctx) error
}

func SetupMiddlewares(app *fiber.App, queue mq.MessageQueue) {
	app.Use(cors.New())

	app.Use(logger.New(logger.Config{
		Format:     "${time} ${status} | ${method} ${path} | ${error} \n",
		TimeFormat: "2 Jan 2006 15:04:05",
		TimeZone:   "local",
	}))

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("messageQueue", queue)
		return c.Next()
	})
}

func GetOptionalMiddlewares(db *gorm.DB) OptionalMiddlewares {
	return OptionalMiddlewares{
		Transaction: GetTransactionMiddleware(db),
		Auth:        AuthMiddleware,
	}
}
