package service

import (
	"context"
	"errors"
	"github.com/secure-notes/internal/domain"
	"github.com/secure-notes/internal/security"
	"gorm.io/gorm"
	"strings"
	"time"
)

type UserAuth struct {
	repo domain.UserRepository
	jwt  *security.JWTManager
}

var (
	ErrInvalidEmail       = errors.New("invalid email")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrPasswordTooShort   = errors.New("password is short")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func NewUserAuth(repo domain.UserRepository, jwtm *security.JWTManager) *UserAuth {
	return &UserAuth{repo: repo, jwt: jwtm}
}
func (u *UserAuth) Register(ctx context.Context, user domain.User) (domain.User, error) {
	if strings.TrimSpace(user.Email) == "" {
		return domain.User{}, ErrInvalidEmail
	}
	if strings.TrimSpace(user.PasswordHash) == "" {
		return domain.User{}, ErrInvalidPassword
	}
	if len(strings.TrimSpace(user.PasswordHash)) < 16 {
		return domain.User{}, ErrPasswordTooShort
	}
	var err error
	user.PasswordHash, err = security.HashPassword(user.PasswordHash, security.DefaultArgon2Params())
	if err != nil {
		return domain.User{}, ErrInvalidCredentials
	}
	_, err = u.repo.Register(ctx, user)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.User{}, ErrEmailAlreadyExists
		}
		return domain.User{}, err
	}
	return domain.User{}, nil
}

type LoginResult struct {
	AccessToken string
	ExpiresAt   time.Time
}

func (u *UserAuth) Login(ctx context.Context, user domain.User) (LoginResult, error) {
	if strings.TrimSpace(user.Email) == "" {
		return LoginResult{}, ErrInvalidEmail
	}
	if strings.TrimSpace(user.PasswordHash) == "" {
		return LoginResult{}, ErrInvalidPassword
	}
	if len(strings.TrimSpace(user.PasswordHash)) < 16 {
		return LoginResult{}, ErrPasswordTooShort
	}
	dbUser, err := u.repo.GetByEmail(ctx, user.Email)

	status, err := security.VerifyPassword(user.PasswordHash, dbUser.PasswordHash)
	if err != nil || !status {
		return LoginResult{}, ErrInvalidCredentials
	}
	token, exp, err := u.jwt.Sign(dbUser.ID)
	if err != nil {
		return LoginResult{}, err
	}
	return LoginResult{
		AccessToken: token,
		ExpiresAt:   exp,
	}, nil
}
