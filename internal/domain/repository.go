package domain

import "context"

// NoteRepository
type NoteRepository interface {
	Create(ctx context.Context, note Note) (Note, error)
	GetByID(ctx context.Context, id int64) (Note, error)
	Update(ctx context.Context, note Note) error
	RemoveByID(ctx context.Context, id int64) error
}

// UserRepository
type UserRepository interface {
	Register(ctx context.Context, user User) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
}
