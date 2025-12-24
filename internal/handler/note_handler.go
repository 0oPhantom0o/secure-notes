package handler

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/secure-notes/internal/domain"
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json body"})
	}

	note := domain.Note{
		Title:   req.Title,
		Content: req.Content,
	}
	created, err := h.svc.CreateNote(c.Context(), note)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(created)
}

func (h NoteHandler) GetByID(c *fiber.Ctx) error {
	noteID := c.Params("id")
	//if strings.TrimSpace(noteID) == "" {
	//	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	//}
	id, err := strconv.Atoi(noteID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	note, err := h.svc.GetByID(c.Context(), int64(id))
	if err != nil {
		if errors.Is(err, domain.ErrNoteNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "note not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(note)
}
