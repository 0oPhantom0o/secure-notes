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
	uidAny := c.Locals(middleware.LocalUserIDKey)
	uid, ok := uidAny.(int64)
	if !ok || uid <= 0 {
		return c.Status(fiber.StatusUnauthorized).
			JSON(response.NewError(response.CodeUnauthorized, "unauthorized"))
	}
	note.UserID = uid
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
	id := c.Params("id")

	noteID, err := idValidator(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewError(response.CodeInvalidID, "invalid id"))
	}
	uidAny := c.Locals(middleware.LocalUserIDKey)
	uid, ok := uidAny.(int64)
	if !ok || uid <= 0 {
		return c.Status(fiber.StatusUnauthorized).
			JSON(response.NewError(response.CodeUnauthorized, "unauthorized"))
	}
	if err = h.svc.RemoveNote(c.Context(), noteID, uid); err != nil {
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
	uidAny := c.Locals(middleware.LocalUserIDKey)
	uid, ok := uidAny.(int64)
	if !ok || uid <= 0 {
		return c.Status(fiber.StatusUnauthorized).
			JSON(response.NewError(response.CodeUnauthorized, "unauthorized"))
	}
	note, err := h.svc.GetByID(c.Context(), id, uid)
	if err != nil {
		if errors.Is(err, domain.ErrNoteNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(response.NewError(response.CodeNoteNotFound, "note not found"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewError(response.CodeInternal, "internal server error"))
	}
	return c.Status(fiber.StatusOK).JSON(note)
}

func (h NoteHandler) List(c *fiber.Ctx) error {
	uidAny := c.Locals("user_id")
	uid, ok := uidAny.(int64)
	if !ok || uid <= 0 {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	notes, err := h.svc.List(c.Context(), uid, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewError(response.CodeInternal, "internal server error"))
	}
	return c.Status(fiber.StatusOK).JSON(notes)
}

func idValidator(id string) (int64, error) {
	return strconv.ParseInt(id, 10, 64)
}
