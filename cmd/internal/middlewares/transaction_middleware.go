package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ctxKey is a type for context keys
type ctxKey string

// txnKey is a context key for database transactions
const TxnKey ctxKey = "db_txn"

func GetTransactionMiddleware(db *gorm.DB) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		tx := db.Begin()
		if tx.Error != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "could not start transaction")
		}

		c.Locals(string(TxnKey), tx)

		// Defer a function to commit or rollback the transaction
		defer func() {
			if r := recover(); r != nil || c.Response().StatusCode() >= 400 {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}()

		return c.Next()
	}
}
