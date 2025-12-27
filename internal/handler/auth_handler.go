package handler

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/secure-notes/internal/domain"
	"github.com/secure-notes/internal/http/response"
	"github.com/secure-notes/internal/service"
	"time"
)

type UserAuthHandler struct {
	svc *service.UserAuth
}
type createUsereAuthReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUserAuthHandler(svc *service.UserAuth) *UserAuthHandler {
	return &UserAuthHandler{svc: svc}
}

func (h UserAuthHandler) Register(c *fiber.Ctx) error {
	var req createUsereAuthReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewError(response.CodeValidation, "bad request"))
	}

	user := domain.User{
		Email:        req.Email,
		PasswordHash: req.Password,
	}
	created, err := h.svc.Register(c.Context(), user)
	if err != nil {
		if errors.Is(err, service.ErrPasswordTooShort) ||
			errors.Is(err, service.ErrInvalidEmail) ||
			errors.Is(err, service.ErrInvalidPassword) {
			return c.Status(fiber.StatusBadRequest).JSON(response.NewError(response.CodeValidation, err.Error()))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewError(response.CodeInternal, "internal server error"))
	}
	return c.Status(fiber.StatusCreated).JSON(created)
}

func (h UserAuthHandler) Login(c *fiber.Ctx) error {
	var req createUsereAuthReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewError(response.CodeValidation, "bad request"))
	}

	user := domain.User{
		Email:        req.Email,
		PasswordHash: req.Password,
	}
	resp, err := h.svc.Login(c.Context(), user)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return c.Status(fiber.StatusBadRequest).JSON(response.NewError(response.CodeUnauthorized, "invalid credentials"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewError(response.CodeInternal, "internal server error"))
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token": resp.AccessToken,
		"token_type":   "Bearer",
		"expires_at":   resp.ExpiresAt.UTC().Format(time.RFC3339),
	})
}
