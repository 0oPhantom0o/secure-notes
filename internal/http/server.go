package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/secure-notes/internal/handler"
	"github.com/secure-notes/internal/service"
)

type NoteApi struct {
	noteSvc *service.NoteService
}

func NewServer(noteHandler *handler.NoteHandler) *fiber.App {
	app := fiber.New()

	api := app.Group("/api")
	api.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
	api.Post("/notes", noteHandler.Create)
	return app
}
