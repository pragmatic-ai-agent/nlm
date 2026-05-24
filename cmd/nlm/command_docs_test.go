package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestCommandReferenceCoversStableCommands(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "docs", "commands.md"))
	if err != nil {
		t.Fatalf("read docs/commands.md: %v", err)
	}

	var missing []string
	for _, cmd := range commandTableEntries() {
		if cmd.surface != surfaceStable || cmd.hidden {
			continue
		}
		want := []byte("`nlm " + cmd.name)
		if !bytes.Contains(data, want) {
			missing = append(missing, cmd.name)
		}
	}
	if len(missing) > 0 {
		t.Fatalf("docs/commands.md missing command table entries: %v", missing)
	}
}
