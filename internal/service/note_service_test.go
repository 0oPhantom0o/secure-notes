package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/secure-notes/internal/domain"
	"github.com/secure-notes/internal/service"
	"gorm.io/gorm"
)

// --- GetByID tests ---

func TestNoteService_GetByID_NotFoundMapsToDomain(t *testing.T) {
	repo := &fakeNoteRepo{
		getByIDFn: func(ctx context.Context, noteID, uid int64) (domain.Note, error) {
			return domain.Note{}, gorm.ErrRecordNotFound
		},
	}
	svc := service.NewNoteService(repo)

	_, err := svc.GetByID(context.Background(), 999, 10) // noteID=999 uid=10
	if !errors.Is(err, domain.ErrNoteNotFound) {
		t.Fatalf("expected ErrNoteNotFound, got: %v", err)
	}
	if repo.lastNoteID != 999 || repo.lastUID != 10 {
		t.Fatalf("ownership args not passed correctly: noteID=%d uid=%d", repo.lastNoteID, repo.lastUID)
	}
}

func TestNoteService_GetByID_Success(t *testing.T) {
	want := domain.Note{ID: 1, UserID: 10, Title: "t", Content: "c"}

	repo := &fakeNoteRepo{
		getByIDFn: func(ctx context.Context, noteID, uid int64) (domain.Note, error) {
			return want, nil
		},
	}
	svc := service.NewNoteService(repo)

	got, err := svc.GetByID(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != want.ID || got.UserID != want.UserID || got.Title != want.Title || got.Content != want.Content {
		t.Fatalf("unexpected note: %+v", got)
	}
}

func TestNoteService_GetByID_OtherErrorPassThrough(t *testing.T) {
	dbErr := errors.New("db down")

	repo := &fakeNoteRepo{
		getByIDFn: func(ctx context.Context, noteID, uid int64) (domain.Note, error) {
			return domain.Note{}, dbErr
		},
	}
	svc := service.NewNoteService(repo)

	_, err := svc.GetByID(context.Background(), 1, 10)
	if !errors.Is(err, dbErr) {
		t.Fatalf("expected db error passthrough, got: %v", err)
	}
}

// --- CreateNote tests --
// --- CreateNote tests ---

func TestNoteService_CreateNote_EmptyTitleReturnsError(t *testing.T) {
	svc := service.NewNoteService(&fakeNoteRepo{})
	_, err := svc.CreateNote(context.Background(), domain.Note{Title: "  ", Content: "content"})

	if !errors.Is(err, service.ErrInvalidTitle) {
		t.Fatalf("expected ErrInvalidTitle, got: %v", err)
	}
}

func TestNoteService_CreateNote_EmptyContentReturnsError(t *testing.T) {
	svc := service.NewNoteService(&fakeNoteRepo{})
	_, err := svc.CreateNote(context.Background(), domain.Note{Title: "title", Content: ""})

	if !errors.Is(err, service.ErrInvalidContent) {
		t.Fatalf("expected ErrInvalidContent, got: %v", err)
	}
}

func TestNoteService_CreateNote_Success(t *testing.T) {
	noteToCreate := domain.Note{Title: "My Note", Content: "Hello world", UserID: 10}
	repo := &fakeNoteRepo{
		createFn: func(ctx context.Context, note domain.Note) (domain.Note, error) {
			note.ID = 500 // فرض می‌کنیم دیتابیس ID تخصیص داده
			return note, nil
		},
	}
	svc := service.NewNoteService(repo)

	created, err := svc.CreateNote(context.Background(), noteToCreate)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created.ID != 500 {
		t.Fatalf("expected ID 500, got %d", created.ID)
	}
	if repo.lastNote.Title != noteToCreate.Title {
		t.Fatalf("wrong note passed to repo")
	}

}

// --- UpdateByID tests ---

func TestNoteService_UpdateByID_NotFound(t *testing.T) {
	repo := &fakeNoteRepo{
		updateFn: func(ctx context.Context, note domain.Note) error {
			return gorm.ErrRecordNotFound
		},
	}
	svc := service.NewNoteService(repo)

	err := svc.UpdateByID(context.Background(), domain.Note{Title: "Update", Content: "New content"})

	if !errors.Is(err, domain.ErrNoteNotFound) {
		t.Fatalf("expected domain.ErrNoteNotFound, got: %v", err)
	}
}

func TestNoteService_UpdateByID_Success(t *testing.T) {
	repo := &fakeNoteRepo{
		updateFn: func(ctx context.Context, note domain.Note) error {
			return nil
		},
	}
	svc := service.NewNoteService(repo)

	updateData := domain.Note{ID: 1, Title: "Updated Title", Content: "Updated Content"}
	err := svc.UpdateByID(context.Background(), updateData)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.lastNote.Title != "Updated Title" {
		t.Fatalf("repo received wrong data")
	} // --- RemoveNote tests ---
} // --- RemoveNote tests ---

func TestNoteService_RemoveNote_NotFound(t *testing.T) {
	repo := &fakeNoteRepo{
		removeByIDFn: func(ctx context.Context, noteID, uid int64) error {
			return gorm.ErrRecordNotFound
		},
	}
	svc := service.NewNoteService(repo)

	err := svc.RemoveNote(context.Background(), 100, 20)

	if !errors.Is(err, domain.ErrNoteNotFound) {
		t.Fatalf("expected ErrNoteNotFound, got: %v", err)
	}
}

func TestNoteService_RemoveNote_Success(t *testing.T) {
	repo := &fakeNoteRepo{
		removeByIDFn: func(ctx context.Context, noteID, uid int64) error {
			return nil
		},
	}
	svc := service.NewNoteService(repo)

	err := svc.RemoveNote(context.Background(), 100, 20)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.lastNoteID != 100 || repo.lastUID != 20 {
		t.Fatalf("correct IDs not passed to repo")
	}
}
