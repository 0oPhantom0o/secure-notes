package service

import (
	"context"
	"errors"

	"github.com/secure-notes/internal/domain"
)

type NoteService struct {
	repo domain.NoteRepository
}

func NewNoteService(repo domain.NoteRepository) *NoteService {
	return &NoteService{repo: repo}
}
func (s *NoteService) CreateNote(note domain.Note) (domain.Note, error) {
	if note.Title == "" {
		return domain.Note{}, errors.New("title is required")
	}
	return s.repo.Create(context.Background(), note)

}
