package service

import (
	"context"
	"errors"
	"gorm.io/gorm"
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
func (s *NoteService) UpdateByID(ctx context.Context, n domain.Note) error {
	if strings.TrimSpace(n.Title) == "" {
		return ErrInvalidTitle
	}
	if strings.TrimSpace(n.Content) == "" {
		return ErrInvalidContent
	}
	err := s.repo.Update(ctx, n)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.ErrNoteNotFound
	}

	return nil
}
func (s *NoteService) RemoveNote(ctx context.Context, id int64) error {
	err := s.repo.RemoveByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNoteNotFound
		}
		return err
	}
	return nil
}
func (s *NoteService) GetByID(ctx context.Context, id int64) (domain.Note, error) {
	note, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Note{}, domain.ErrNoteNotFound
		}
		return domain.Note{}, err
	}
	return note, nil
}
