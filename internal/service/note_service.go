package service

import (
	"context"
	"errors"
	"strings"

	"github.com/secure-notes/internal/domain"
)

type NoteService struct {
	repo domain.NoteRepository
}

var (
	ErrInvalidTitle   = errors.New("title is required")
	ErrInvalidContent = errors.New("content is required")
)

func NewNoteService(repo domain.NoteRepository) *NoteService {
	return &NoteService{repo: repo}
}
func (s *NoteService) CreateNote(ctx context.Context, n domain.Note) (domain.Note, error) {
	if strings.TrimSpace(n.Title) == "" {
		return domain.Note{}, ErrInvalidTitle
	}
	if strings.TrimSpace(n.Content) == "" {
		return domain.Note{}, ErrInvalidContent
	}
	return s.repo.Create(ctx, n)
}
