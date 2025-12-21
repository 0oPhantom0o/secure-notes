package app

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/secure-notes/internal/handler"
	apihttp "github.com/secure-notes/internal/http"
	p "github.com/secure-notes/internal/repository/postgres"
	"github.com/secure-notes/internal/service"
)

func New(ctx context.Context) (*fiber.App, error) {
	db, err := p.NewDB(ctx)
	if err != nil {
		return nil, err
	}
	noteRepo := p.NewNoteRepo(db)
	noteSvc := service.NewNoteService(noteRepo)
	noteHandler := handler.NewHandler(noteSvc)
	app := apihttp.NewServer(noteHandler)
	return app, nil
}
