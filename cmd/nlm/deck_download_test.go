package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestParseDeckDownloadArgs(t *testing.T) {
	opts, notebookID, err := parseDeckDownloadArgs([]string{
		"notebook-123",
		"--id", "artifact-456",
		"--format", "pptx",
		"--output", "deck.pptx",
	})
	if err != nil {
		t.Fatalf("parseDeckDownloadArgs() error = %v", err)
	}
	if notebookID != "notebook-123" {
		t.Fatalf("notebookID = %q, want notebook-123", notebookID)
	}
	if opts.ArtifactID != "artifact-456" || opts.Format != "pptx" || opts.Output != "deck.pptx" {
		t.Fatalf("opts = %+v", opts)
	}
}

func TestParseDeckDownloadArgsRejectsBadFormat(t *testing.T) {
	_, _, err := parseDeckDownloadArgs([]string{
		"notebook-123",
		"--id", "artifact-456",
		"--format", "keynote",
	})
	if err == nil || !strings.Contains(err.Error(), "unsupported format") {
		t.Fatalf("parseDeckDownloadArgs() error = %v, want unsupported format", err)
	}
}

func TestDeckDownloadFallbackPrintsBrowserURL(t *testing.T) {
	tmpHome, err := os.MkdirTemp("", "nlm-test-home-*")
	if err != nil {
		t.Fatalf("failed to create temp home: %v", err)
	}
	defer os.RemoveAll(tmpHome)

	cmd := exec.Command("./nlm_test", "deck", "download", "notebook-123", "--id", "artifact-456", "--format", "pptx")
	cmd.Env = []string{"PATH=" + os.Getenv("PATH"), "HOME=" + tmpHome}
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("deck download unexpectedly succeeded\n%s", output)
	}
	out := string(output)
	if !strings.Contains(out, "https://notebooklm.google.com/notebook/notebook-123") {
		t.Fatalf("deck download did not print browser URL\n%s", out)
	}
	if strings.Contains(out, "Authentication required") {
		t.Fatalf("deck download fallback should not require auth\n%s", out)
	}
}

func TestLegacySlideDeckDownloadPrintsBrowserURL(t *testing.T) {
	tmpHome, err := os.MkdirTemp("", "nlm-test-home-*")
	if err != nil {
		t.Fatalf("failed to create temp home: %v", err)
	}
	defer os.RemoveAll(tmpHome)

	cmd := exec.Command("./nlm_test", "download", "slide-deck", "notebook-123", "--id", "artifact-456", "--format", "pptx", "--output", "deck.pptx")
	cmd.Env = []string{"PATH=" + os.Getenv("PATH"), "HOME=" + tmpHome}
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("legacy download unexpectedly succeeded\n%s", output)
	}
	out := string(output)
	if !strings.Contains(out, "nlm: 'download slide-deck' is deprecated; use 'deck download'") {
		t.Fatalf("legacy download did not print compatibility warning\n%s", out)
	}
	if !strings.Contains(out, "https://notebooklm.google.com/notebook/notebook-123") {
		t.Fatalf("legacy download did not print browser URL\n%s", out)
	}
}
