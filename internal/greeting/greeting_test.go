package greeting

import "testing"

func TestBuild(t *testing.T) {
	got, err := Build("Mihir")
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	want := "Hello, Mihir!"
	if got != want {
		t.Fatalf("Build() = %q, want %q", got, want)
	}
}

func TestBuildTrimsName(t *testing.T) {
	got, err := Build("  Codex  ")
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	want := "Hello, Codex!"
	if got != want {
		t.Fatalf("Build() = %q, want %q", got, want)
	}
}

func TestBuildRejectsEmptyName(t *testing.T) {
	if _, err := Build("   "); err == nil {
		t.Fatal("Build returned nil error for empty name")
	}
}
