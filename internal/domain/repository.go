package domain

import "context"

// NoteRepository
type NoteRepository interface {
	Create(ctx context.Context, note Note) (Note, error)
	List(ctx context.Context, uid int64, limit, offset int) ([]Note, error)
	Update(ctx context.Context, note Note) error
	GetByID(ctx context.Context, noteID, uid int64) (Note, error)
	RemoveByID(ctx context.Context, noteID, uid int64) error
}

// UserRepository
type UserRepository interface {
	Register(ctx context.Context, user User) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
}
