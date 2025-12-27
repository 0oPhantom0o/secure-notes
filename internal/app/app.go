package app

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/secure-notes/internal/security"
	"os"
	"time"

	"github.com/secure-notes/internal/handler"
	apihttp "github.com/secure-notes/internal/http"
	p "github.com/secure-notes/internal/repository/postgres"
	"github.com/secure-notes/internal/service"
)

func New(ctx context.Context) (*fiber.App, error) {
	jwtm := security.NewJWTManager(os.Getenv("JWT_SECRET"), "secure-notes", time.Hour)
	db, err := p.NewDB(ctx)
	if err != nil {
		return nil, err
	}
	noteRepo := p.NewNoteRepo(db)
	userRepo := p.NewUserRepo(db)
	noteSvc := service.NewNoteService(noteRepo)
	userSvc := service.NewUserAuth(userRepo, jwtm)
	noteHandler := handler.NewHandler(noteSvc)
	userHandler := handler.NewUserAuthHandler(userSvc)
	app := apihttp.NewServer(noteHandler, userHandler, jwtm)
	return app, nil
}
