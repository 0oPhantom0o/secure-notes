package handler

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/secure-notes/internal/domain"
	"github.com/secure-notes/internal/http/middleware"
	"github.com/secure-notes/internal/http/response"
	"github.com/secure-notes/internal/service"
	"strconv"
)

type NoteHandler struct {
	svc *service.NoteService
}
type createNoteReq struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func NewHandler(svc *service.NoteService) *NoteHandler {
	return &NoteHandler{svc: svc}
}

func (h NoteHandler) Create(c *fiber.Ctx) error {
	var req createNoteReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewError(response.CodeValidation, "bad request"))
	}

	note := domain.Note{
		Title:   req.Title,
		Content: req.Content,
	}
	uidAny := c.Locals(middleware.LocalUserIDKey)
	uid, ok := uidAny.(int64)
	if !ok || uid <= 0 {
		return c.Status(fiber.StatusUnauthorized).
			JSON(response.NewError(response.CodeUnauthorized, "unauthorized"))
	}
	note.UserID = uid
	created, err := h.svc.CreateNote(c.Context(), note)
	if err != nil {
		if errors.Is(err, service.ErrInvalidContent) || errors.Is(err, service.ErrInvalidTitle) {
			return c.Status(fiber.StatusBadRequest).JSON(response.NewError(response.CodeValidation, err.Error()))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewError(response.CodeInternal, "internal server error"))
	}
	return c.Status(fiber.StatusCreated).JSON(created)
}

func (h NoteHandler) UpdateByID(c *fiber.Ctx) error {
	noteID := c.Params("id")

	id, err := idValidator(noteID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewError(response.CodeInvalidID, "invalid id"))
	}
	var req createNoteReq
	if err = c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewError(response.CodeValidation, "bad request"))
	}

	note := domain.Note{
		ID:      id,
		Title:   req.Title,
		Content: req.Content,
	}
	err = h.svc.UpdateByID(c.Context(), note)
	if err != nil {
		if errors.Is(err, domain.ErrNoteNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(response.NewError(response.CodeNoteNotFound, "note not found"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewError(response.CodeInternal, "internal server error"))
	}
	return c.SendStatus(fiber.StatusNoContent)
}
func (h NoteHandler) DeleteByID(c *fiber.Ctx) error {
	noteID := c.Params("id")

	id, err := idValidator(noteID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewError(response.CodeInvalidID, "invalid id"))
	}
	err = h.svc.RemoveNote(c.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNoteNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(response.NewError(response.CodeNoteNotFound, "note not found"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewError(response.CodeInternal, "internal server error"))
	}
	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"status": true})

}

func (h NoteHandler) GetByID(c *fiber.Ctx) error {
	noteID := c.Params("id")

	id, err := idValidator(noteID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewError(response.CodeInvalidID, "invalid id"))
	}
	note, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNoteNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(response.NewError(response.CodeNoteNotFound, "note not found"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewError(response.CodeInternal, "internal server error"))
	}
	return c.Status(fiber.StatusOK).JSON(note)
}
func idValidator(id string) (int64, error) {
	return strconv.ParseInt(id, 10, 64)
}
