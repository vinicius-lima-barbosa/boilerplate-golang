package middleware

import (
	"github.com/gofiber/fiber/v2"
	config "github.com/vinicius-lima-barbosa/boilerplate-golang/internal/database"
)

const DBTransactionKey = "db_trx"

func DBTransactionMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		txHandle := config.DB.DB.Begin()
		c.Locals(DBTransactionKey, txHandle)

		defer func() {
			if r := recover(); r != nil {
				txHandle.Rollback()
				panic(r)
			}
		}()

		err := c.Next()
		if err != nil {
			txHandle.Rollback()
			return err
		}

		if c.Response().StatusCode() >= fiber.StatusBadRequest {
			txHandle.Rollback()
			return nil
		}

		if err := txHandle.Commit().Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to commit database transaction",
			})
		}

		return nil
	}
}
