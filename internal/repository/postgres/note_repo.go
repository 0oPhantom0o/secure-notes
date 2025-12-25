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
	if err := r.db.WithContext(ctx).First(&note, "id = ?", id).Error; err != nil {
		return domain.Note{}, err
	}
	return note, nil
}
func (r NoteRepo) Update(ctx context.Context, note domain.Note) error {
	tx := r.db.WithContext(ctx).
		Model(&domain.Note{}).
		Where("id = ?", note.ID).
		Updates(map[string]any{
			"title":   note.Title,
			"content": note.Content,
		})

	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r NoteRepo) RemoveByID(ctx context.Context, id int64) error {
	var note domain.Note
	Rerr := r.db.WithContext(ctx).Delete(&note, id)
	if Rerr.Error != nil {
		return Rerr.Error
	}
	if Rerr.RowsAffected <= 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
