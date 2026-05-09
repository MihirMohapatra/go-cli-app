package todos

import "testing"

func TestStoreCRUD(t *testing.T) {
	store := NewStore()

	created, err := store.Create("Learn Go")
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if created.ID != 1 {
		t.Fatalf("Create ID = %d, want 1", created.ID)
	}

	got, ok := store.Get(created.ID)
	if !ok {
		t.Fatal("Get could not find created todo")
	}
	if got.Title != "Learn Go" || got.Completed {
		t.Fatalf("Get = %+v, want title Learn Go and completed false", got)
	}

	updated, err := store.Update(created.ID, "Ship API", true)
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if updated.Title != "Ship API" || !updated.Completed {
		t.Fatalf("Update = %+v, want title Ship API and completed true", updated)
	}

	if !store.Delete(created.ID) {
		t.Fatal("Delete returned false for existing todo")
	}
	if _, ok := store.Get(created.ID); ok {
		t.Fatal("Get found deleted todo")
	}
}

func TestStoreRejectsEmptyTitle(t *testing.T) {
	store := NewStore()

	if _, err := store.Create("   "); err == nil {
		t.Fatal("Create returned nil error for empty title")
	}
}
