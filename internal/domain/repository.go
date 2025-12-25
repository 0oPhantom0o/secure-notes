package domain

import "context"

type NoteRepository interface {
	Create(ctx context.Context, note Note) (Note, error)
	GetByID(ctx context.Context, id int64) (Note, error)
	Update(ctx context.Context, note Note) error
	RemoveByID(ctx context.Context, id int64) error
}
