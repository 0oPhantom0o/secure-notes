package service_test

import (
	"context"
	"github.com/secure-notes/internal/domain"
)

type fakeNoteRepo struct {
	getByIDFn    func(ctx context.Context, noteID, uid int64) (domain.Note, error)
	updateFn     func(ctx context.Context, note domain.Note) error
	removeByIDFn func(ctx context.Context, noteID, uid int64) error
	createFn     func(ctx context.Context, note domain.Note) (domain.Note, error)

	lastNoteID int64
	lastUID    int64
	lastNote   domain.Note
}

func (f *fakeNoteRepo) GetByID(ctx context.Context, noteID, uid int64) (domain.Note, error) {
	if f.getByIDFn == nil {
		panic("fakeNoteRepo.getByIDFn not set")
	}
	f.lastNoteID, f.lastUID = noteID, uid
	return f.getByIDFn(ctx, noteID, uid)
}

func (f *fakeNoteRepo) Update(ctx context.Context, note domain.Note) error {
	if f.updateFn == nil {
		panic("fakeNoteRepo.updateFn not set")
	}
	f.lastNote = note
	return f.updateFn(ctx, note)
}

func (f *fakeNoteRepo) RemoveByID(ctx context.Context, noteID, uid int64) error {
	if f.removeByIDFn == nil {
		panic("fakeNoteRepo.removeByIDFn not set")
	}
	f.lastNoteID, f.lastUID = noteID, uid
	return f.removeByIDFn(ctx, noteID, uid)
}

func (f *fakeNoteRepo) Create(ctx context.Context, note domain.Note) (domain.Note, error) {
	if f.createFn == nil {
		panic("fakeNoteRepo.createFn not set")
	}
	f.lastNote = note
	return f.createFn(ctx, note)
}
