package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/secure-notes/internal/handler"
	"github.com/secure-notes/internal/http/middleware"
	"github.com/secure-notes/internal/security"
	"github.com/secure-notes/internal/service"
)

type NoteApi struct {
	noteSvc *service.NoteService
}

const (
	prefix   = "/api/v1"
	pingPath = "/healthz"
	pathNote = "/notes"
	pathAuth = "/auth"
	register = "/register"
	login    = "/login"
)

func NewServer(noteHandler *handler.NoteHandler, userHandler *handler.UserAuthHandler, jwtm *security.JWTManager) *fiber.App {
	app := fiber.New()
	api := app.Group(prefix)
	api.Get(pingPath, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	notes := api.Group(pathNote, middleware.AuthRequired(jwtm))
	notes.Post("/", noteHandler.Create)
	notes.Get(pathNote+"/:id", noteHandler.GetByID)
	notes.Put(pathNote+"/:id", noteHandler.UpdateByID)
	notes.Delete(pathNote+"/:id", noteHandler.DeleteByID)

	auth := api.Group(pathAuth)
	auth.Post(register, userHandler.Register)
	auth.Post(login, userHandler.Login)

	return app
}
