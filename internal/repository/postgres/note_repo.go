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

func (r NoteRepo) List(ctx context.Context, uid int64, limit, offset int) ([]domain.Note, error) {
	var notes []domain.Note
	if err := r.db.WithContext(ctx).Where("user_id = ?", uid).Limit(limit).Offset(offset).Order("id DESC").
		Find(&notes).Error; err != nil {

		return nil, err
	}
	return notes, nil
}

func (r NoteRepo) Create(ctx context.Context, note domain.Note) (domain.Note, error) {
	if err := r.db.WithContext(ctx).Create(&note).Error; err != nil {
		return domain.Note{}, err
	}
	return note, nil
}

func (r NoteRepo) GetByID(ctx context.Context, noteID, uid int64) (domain.Note, error) {
	var note domain.Note
	if err := r.db.WithContext(ctx).First(&note, "id = ? and user_id = ?", noteID, uid).Error; err != nil {
		return domain.Note{}, err
	}
	return note, nil
}
func (r NoteRepo) Update(ctx context.Context, note domain.Note) error {
	tx := r.db.WithContext(ctx).
		Model(&domain.Note{}).
		Where("id = ? and user_id = ? ", note.ID, note.UserID).
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

func (r NoteRepo) RemoveByID(ctx context.Context, noteID, uid int64) error {
	var note domain.Note
	Rerr := r.db.WithContext(ctx).Delete(&note).Where("id = ? and user_id = ? ", noteID, uid)
	if Rerr.Error != nil {
		return Rerr.Error
	}
	if Rerr.RowsAffected <= 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
