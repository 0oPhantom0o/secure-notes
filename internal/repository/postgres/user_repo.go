package postgres

import (
	"context"
	"errors"
	"github.com/secure-notes/internal/domain"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r UserRepo) Register(ctx context.Context, user domain.User) (domain.User, error) {
	if err := r.db.WithContext(ctx).Create(&user).Error; err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (r UserRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, gorm.ErrRecordNotFound
		}
		return domain.User{}, err
	}
	return user, nil
}
