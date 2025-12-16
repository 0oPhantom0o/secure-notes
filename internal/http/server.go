package http

import "github.com/gofiber/fiber/v2"

func NewServer() *fiber.App {
	app := fiber.New()

	api := app.Group("/api")
	api.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	return app
}
