package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/secure-notes/internal/handler"
	"github.com/secure-notes/internal/service"
)

type NoteApi struct {
	noteSvc *service.NoteService
}

const prefix = "/api/v1"
const pingPath = "/healthz"
const pathNote = "/notes"

func NewServer(noteHandler *handler.NoteHandler) *fiber.App {
	app := fiber.New()

	api := app.Group(prefix)
	api.Get(pingPath, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
	api.Post(pathNote, noteHandler.Create)
	return app
}
