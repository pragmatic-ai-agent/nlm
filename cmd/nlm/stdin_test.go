package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestReadLinesFromReader(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []string
	}{
		{"empty", "", nil},
		{"single line", "abc\n", []string{"abc"}},
		{"trailing newline omitted", "abc", []string{"abc"}},
		{"multiple lines", "a\nb\nc\n", []string{"a", "b", "c"}},
		{"blank lines skipped", "a\n\n\nb\n", []string{"a", "b"}},
		{"comments skipped", "a\n# comment\nb\n", []string{"a", "b"}},
		{"leading whitespace trimmed", "  a\n\tb\n", []string{"a", "b"}},
		{"column splits on first whitespace", "id1 Title One\nid2\tType\n", []string{"id1", "id2"}},
		{"mixed with blanks and comments", "\n# header\nuuid-1\n\n  # more\nuuid-2  extra\n", []string{"uuid-1", "uuid-2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readLinesFromReader(strings.NewReader(tt.in))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestResolveIDList(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want []string
	}{
		{"empty", "", nil},
		{"single id", "abc", []string{"abc"}},
		{"comma list", "a,b,c", []string{"a", "b", "c"}},
		{"comma list with spaces", "a, b , c", []string{"a", "b", "c"}},
		{"comma list trims empties", "a,,b,", []string{"a", "b"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveIDList(tt.arg)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}

// Note: resolveIDList("-") requires non-TTY stdin; exercised at integration level.

type fakeTTY struct {
	r *strings.Reader
	w bytes.Buffer
}

func (f *fakeTTY) Read(p []byte) (int, error)  { return f.r.Read(p) }
func (f *fakeTTY) Write(p []byte) (int, error) { return f.w.Write(p) }
func (f *fakeTTY) Close() error                { return nil }

func TestConfirmActionUsesControllingTTY(t *testing.T) {
	oldYes := yes
	yes = false
	t.Cleanup(func() { yes = oldYes })

	oldOpen := openControllingTTY
	tty := &fakeTTY{r: strings.NewReader("y\n")}
	openControllingTTY = func() (io.ReadWriteCloser, error) { return tty, nil }
	t.Cleanup(func() { openControllingTTY = oldOpen })

	oldStdin := os.Stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.WriteString("pipeline-data\n"); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	os.Stdin = r
	t.Cleanup(func() {
		os.Stdin = oldStdin
		r.Close()
	})

	if !confirmAction("delete?") {
		t.Fatal("confirmAction returned false; want true")
	}
	if got := tty.w.String(); !strings.Contains(got, "delete? [y/N]") {
		t.Fatalf("tty prompt = %q, want confirmation prompt", got)
	}
	left, err := io.ReadAll(os.Stdin)
	if err != nil {
		t.Fatal(err)
	}
	if string(left) != "pipeline-data\n" {
		t.Fatalf("stdin was consumed: got %q", left)
	}
}

func TestConfirmActionNoTTYLeavesStdin(t *testing.T) {
	oldYes := yes
	yes = false
	t.Cleanup(func() { yes = oldYes })

	oldOpen := openControllingTTY
	openControllingTTY = func() (io.ReadWriteCloser, error) {
		return nil, errors.New("no tty")
	}
	t.Cleanup(func() { openControllingTTY = oldOpen })

	oldStdin := os.Stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.WriteString("y\n"); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	os.Stdin = r
	t.Cleanup(func() {
		os.Stdin = oldStdin
		r.Close()
	})

	if confirmAction("delete?") {
		t.Fatal("confirmAction returned true without a tty")
	}
	left, err := io.ReadAll(os.Stdin)
	if err != nil {
		t.Fatal(err)
	}
	if string(left) != "y\n" {
		t.Fatalf("stdin was consumed: got %q", left)
	}
}
