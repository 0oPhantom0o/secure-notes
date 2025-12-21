package postgres

import (
	"context"

	"github.com/secure-notes/internal/domain"
	"gorm.io/gorm"
)

type NoteRepo struct {
	db *gorm.DB
}

func NewNoteRepo(db *gorm.DB) *NoteRepo {
	return &NoteRepo{db: db}
}

func (r NoteRepo) Create(ctx context.Context, note domain.Note) (domain.Note, error) {
	if err := r.db.WithContext(ctx).Create(&note).Error; err != nil {
		return domain.Note{}, err
	}
	return note, nil
}
func (r NoteRepo) GetByID(ctx context.Context, id int64) (domain.Note, error) {
	var note domain.Note
	if err := r.db.WithContext(ctx).First(note, "id = ?", id).Error; err != nil {
		return domain.Note{}, err
	}
	return note, nil
}

//yadegari from amir
