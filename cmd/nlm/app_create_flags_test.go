package main

import (
	"testing"

	"github.com/tmc/nlm/internal/notebooklm/api"
)

func TestParseAppCreateArgsWithOptions(t *testing.T) {
	t.Parallel()

	opts, positional, err := parseAppCreateArgsWithOptions([]string{
		"--type", "mindmap",
		"--instructions", "focus on architecture",
		"--source-ids", "src-1,src-2",
		"nb-1",
	}, globalOptions{})
	if err != nil {
		t.Fatalf("parseAppCreateArgsWithOptions: %v", err)
	}
	if opts.Type != "mindmap" || opts.Instructions != "focus on architecture" {
		t.Fatalf("type/instructions = %q/%q", opts.Type, opts.Instructions)
	}
	if opts.Selectors.SourceIDs != "src-1,src-2" {
		t.Fatalf("source ids = %q, want src-1,src-2", opts.Selectors.SourceIDs)
	}
	if len(positional) != 1 || positional[0] != "nb-1" {
		t.Fatalf("positional = %v, want [nb-1]", positional)
	}
}

func TestParseAppCreateArgsUsesPositionalInstructions(t *testing.T) {
	t.Parallel()

	opts, positional, err := parseAppCreateArgsWithOptions([]string{
		"--type", "prototype",
		"nb-1",
		"build", "a", "study", "app",
	}, globalOptions{})
	if err != nil {
		t.Fatalf("parseAppCreateArgsWithOptions: %v", err)
	}
	if opts.Instructions != "build a study app" {
		t.Fatalf("instructions = %q, want positional join", opts.Instructions)
	}
	if len(positional) != 1 || positional[0] != "nb-1" {
		t.Fatalf("positional = %v, want [nb-1]", positional)
	}
}

func TestParseSlideDeckFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in      string
		want    api.SlideDeckFormat
		wantErr bool
	}{
		{"", api.SlideDeckFormatDetailed, false},
		{"detailed", api.SlideDeckFormatDetailed, false},
		{"DETAILED", api.SlideDeckFormatDetailed, false},
		{"detail", api.SlideDeckFormatDetailed, false},
		{"presenter", api.SlideDeckFormatPresenter, false},
		{" Presenter ", api.SlideDeckFormatPresenter, false},
		{"sparse", api.SlideDeckFormatPresenter, false},
		{"bogus", 0, true},
	}
	for _, tt := range tests {
		got, err := parseSlideDeckFormat(tt.in)
		if tt.wantErr {
			if err == nil {
				t.Errorf("parseSlideDeckFormat(%q) = %v, want error", tt.in, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("parseSlideDeckFormat(%q): %v", tt.in, err)
			continue
		}
		if got != tt.want {
			t.Errorf("parseSlideDeckFormat(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestParseSlidesCreateArgs(t *testing.T) {
	t.Parallel()

	// Flags before and after the positional notebook id; instructions optional.
	opts, positional, err := parseSlidesCreateArgs([]string{
		"--format", "presenter",
		"nb-1",
		"focus", "on", "the", "results",
		"--source-match", "^spec/",
	}, globalOptions{})
	if err != nil {
		t.Fatalf("parseSlidesCreateArgs: %v", err)
	}
	if opts.Format != "presenter" {
		t.Fatalf("format = %q, want presenter", opts.Format)
	}
	if opts.Selectors.SourceMatch != "^spec/" {
		t.Fatalf("source-match = %q, want ^spec/", opts.Selectors.SourceMatch)
	}
	if len(positional) != 5 || positional[0] != "nb-1" {
		t.Fatalf("positional = %v, want [nb-1 focus on the results]", positional)
	}

	// Notebook id alone (no instructions, no format) is valid.
	if _, _, err := parseSlidesCreateArgs([]string{"nb-1"}, globalOptions{}); err != nil {
		t.Fatalf("parseSlidesCreateArgs(nb only): %v", err)
	}

	// Missing notebook id is an error.
	if _, _, err := parseSlidesCreateArgs([]string{"--format", "detailed"}, globalOptions{}); err == nil {
		t.Fatal("parseSlidesCreateArgs with no notebook id: want error")
	}

	// Invalid format is rejected at parse time.
	if _, _, err := parseSlidesCreateArgs([]string{"--format", "bogus", "nb-1"}, globalOptions{}); err == nil {
		t.Fatal("parseSlidesCreateArgs with bad format: want error")
	}
}

func TestParseAudioVideoOptions(t *testing.T) {
	t.Parallel()

	aopts, apos, err := parseAudioCreateArgs([]string{
		"--length", "long",
		"--language", "es",
		"--audio-type", "debate",
		"nb-1",
		"compare the sources",
	})
	if err != nil {
		t.Fatalf("parseAudioCreateArgs: %v", err)
	}
	if aopts.Length != "long" || aopts.Language != "es" || aopts.AudioType != "debate" {
		t.Fatalf("audio opts = %+v", aopts)
	}
	if len(apos) != 2 {
		t.Fatalf("audio positional = %v", apos)
	}

	vopts, vpos, err := parseVideoCreateArgs([]string{
		"--style", "whiteboard",
		"--language", "fr",
		"nb-1",
		"explain visually",
	})
	if err != nil {
		t.Fatalf("parseVideoCreateArgs: %v", err)
	}
	if vopts.Style != "whiteboard" || vopts.Language != "fr" {
		t.Fatalf("video opts = %+v", vopts)
	}
	if len(vpos) != 2 {
		t.Fatalf("video positional = %v", vpos)
	}
}
