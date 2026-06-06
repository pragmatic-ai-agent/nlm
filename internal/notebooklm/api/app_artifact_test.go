package api

import "testing"

func TestParseAppArtifactKind(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want AppArtifactKind
	}{
		{"prototype", AppArtifactKindPrototype},
		{"notebook-app", AppArtifactKindPrototype},
		{"mindmap", AppArtifactKindMindmap},
		{"mind-map", AppArtifactKindMindmap},
		{"canvas", AppArtifactKindCanvas},
	}
	for _, tt := range tests {
		got, err := ParseAppArtifactKind(tt.in)
		if err != nil {
			t.Fatalf("ParseAppArtifactKind(%q): %v", tt.in, err)
		}
		if got != tt.want {
			t.Fatalf("ParseAppArtifactKind(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestParseAppArtifactKindRejectsUnknown(t *testing.T) {
	t.Parallel()

	if _, err := ParseAppArtifactKind("flashcards"); err == nil {
		t.Fatal("ParseAppArtifactKind(flashcards) succeeded, want error")
	}
}

func TestParseCreatedArtifactID(t *testing.T) {
	t.Parallel()

	got, err := parseCreatedArtifactID([]byte(`[[ "artifact-1", "Title", 5 ]]`))
	if err != nil {
		t.Fatalf("parseCreatedArtifactID: %v", err)
	}
	if got != "artifact-1" {
		t.Fatalf("artifact id = %q, want artifact-1", got)
	}
}

// TestParseCreatedArtifactIDRejectsEmpty verifies that a blank id — what the
// server returns when a create is rejected without an RPC-level error (e.g.
// quota exhausted) — surfaces as an error instead of a silent empty id.
func TestParseCreatedArtifactIDRejectsEmpty(t *testing.T) {
	t.Parallel()

	for _, resp := range []string{`[""]`, `[["", "Title", 5]]`, `[]`, `[[]]`} {
		got, err := parseCreatedArtifactID([]byte(resp))
		if err == nil {
			t.Errorf("parseCreatedArtifactID(%s) = %q, nil; want error", resp, got)
		}
		if got != "" {
			t.Errorf("parseCreatedArtifactID(%s) id = %q, want empty", resp, got)
		}
	}
}
