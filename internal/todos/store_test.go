package todos

import (
	"context"
	"testing"
)

func TestStoreCRUD(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore()

	created, err := store.Create(ctx, "Learn Go")
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if created.ID != 1 {
		t.Fatalf("Create ID = %d, want 1", created.ID)
	}

	got, err := store.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if got.Title != "Learn Go" || got.Completed {
		t.Fatalf("Get = %+v, want title Learn Go and completed false", got)
	}

	updated, err := store.Update(ctx, created.ID, "Ship API", true)
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if updated.Title != "Ship API" || !updated.Completed {
		t.Fatalf("Update = %+v, want title Ship API and completed true", updated)
	}

	if err := store.Delete(ctx, created.ID); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	if _, err := store.Get(ctx, created.ID); err == nil {
		t.Fatal("Get returned nil error for deleted todo")
	}
}

func TestStoreRejectsEmptyTitle(t *testing.T) {
	store := NewMemoryStore()

	if _, err := store.Create(context.Background(), "   "); err == nil {
		t.Fatal("Create returned nil error for empty title")
	}
}
