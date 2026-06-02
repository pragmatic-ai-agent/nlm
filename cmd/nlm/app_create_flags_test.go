package main

import "testing"

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
