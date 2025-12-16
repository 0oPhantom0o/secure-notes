package postgres

import (
	"context"

	"github.com/secure-notes/internal/domain"
	"gorm.io/gorm"
)

type NoteRepository struct {
	db *gorm.DB
}

func (r NoteRepository) Create(ctx context.Context, note domain.Note) (domain.Note, error) {
	return domain.Note{}, nil
}
func (r NoteRepository) GetByID(ctx context.Context, id string) (domain.Note, error) {
	return domain.Note{}, nil
}
