package server

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
)

func LoggerMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		log.Default().Printf("[%s] %s %d %s (%v)",
			c.Method(),
			c.OriginalURL(),
			c.Response().StatusCode(),
			string(c.Response().Body()),
			time.Since(start))

		return err
	}
}
