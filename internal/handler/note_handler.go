package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/secure-notes/internal/domain"
	"github.com/secure-notes/internal/service"
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
