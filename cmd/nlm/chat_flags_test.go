package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseSourceSelectionArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantPos []string
		wantSel selectorOptions
		wantErr string
	}{
		{
			name:    "selectors without positional source ids",
			args:    []string{"nb", "--source-match", "^spec/"},
			wantPos: []string{"nb"},
			wantSel: selectorOptions{SourceMatch: "^spec/"},
		},
		{
			name:    "positional source ids still work",
			args:    []string{"nb", "src-1", "src-2"},
			wantPos: []string{"nb", "src-1", "src-2"},
		},
		{
			name:    "missing source ids and selectors",
			args:    []string{"nb"},
			wantErr: "missing source ids or selectors",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotPos, err := parseSourceSelectionArgsWithOptions(tt.args, globalOptions{})
			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Fatalf("parseSourceSelectionArgs(%q) error = %v, want %q", tt.args, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseSourceSelectionArgs(%q) error = %v", tt.args, err)
			}
			if got.Selectors != tt.wantSel {
				t.Fatalf("parseSourceSelectionArgs(%q) selectors = %+v, want %+v", tt.args, got.Selectors, tt.wantSel)
			}
			if len(gotPos) != len(tt.wantPos) {
				t.Fatalf("parseSourceSelectionArgs(%q) positional = %q, want %q", tt.args, gotPos, tt.wantPos)
			}
			for i := range gotPos {
				if gotPos[i] != tt.wantPos[i] {
					t.Fatalf("parseSourceSelectionArgs(%q) positional = %q, want %q", tt.args, gotPos, tt.wantPos)
				}
			}
		})
	}
}

func TestParseGenerateChatArgs(t *testing.T) {
	got, gotPos, err := parseGenerateChatArgsWithOptions([]string{
		"nb",
		"why",
		"--conversation", "conv-1",
		"--thinking",
		"--source-match", "^spec/",
		"now",
	}, globalOptions{})
	if err != nil {
		t.Fatalf("parseGenerateChatArgs error = %v", err)
	}
	if got.ConversationID != "conv-1" || !got.Render.ShowThinking || got.Selectors.SourceMatch != "^spec/" {
		t.Fatalf("parseGenerateChatArgs opts = %+v", got)
	}
	wantPos := []string{"nb", "why", "now"}
	if len(gotPos) != len(wantPos) {
		t.Fatalf("parseGenerateChatArgs positional = %q, want %q", gotPos, wantPos)
	}
	for i := range gotPos {
		if gotPos[i] != wantPos[i] {
			t.Fatalf("parseGenerateChatArgs positional = %q, want %q", gotPos, wantPos)
		}
	}
}

func TestParseGenerateChatArgsPromptFile(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "long flag",
			args: []string{"--prompt-file", "prompt.txt", "nb"},
			want: "prompt.txt",
		},
		{
			name: "short flag",
			args: []string{"nb", "-f", "-"},
			want: "-",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotPos, err := parseGenerateChatArgsWithOptions(tt.args, globalOptions{})
			if err != nil {
				t.Fatalf("parseGenerateChatArgs error = %v", err)
			}
			if got.PromptFile != tt.want {
				t.Fatalf("prompt file = %q, want %q", got.PromptFile, tt.want)
			}
			if len(gotPos) != 1 || gotPos[0] != "nb" {
				t.Fatalf("positional = %q, want [nb]", gotPos)
			}
		})
	}
}

func TestParseChatArgs(t *testing.T) {
	got, gotPos, err := parseChatArgsWithOptions([]string{
		"nb",
		"--prompt-file", "prompt.txt",
		"--history",
		"--citations", "tail",
		"--source-ids", "a,b",
	}, globalOptions{})
	if err != nil {
		t.Fatalf("parseChatArgs error = %v", err)
	}
	if got.PromptFile != "prompt.txt" || !got.ShowHistory || got.Render.CitationMode != "tail" || got.Selectors.SourceIDs != "a,b" {
		t.Fatalf("parseChatArgs opts = %+v", got)
	}
	if len(gotPos) != 1 || gotPos[0] != "nb" {
		t.Fatalf("parseChatArgs positional = %q, want [nb]", gotPos)
	}
}

func TestParseGenerateReportArgs(t *testing.T) {
	got, gotPos, err := parseGenerateReportArgsWithOptions([]string{
		"nb",
		"--sections", "3",
		"--prompt", "# {topic}",
		"--source-match", "^guide/",
	}, globalOptions{})
	if err != nil {
		t.Fatalf("parseGenerateReportArgs error = %v", err)
	}
	if got.Sections != 3 || got.Prompt != "# {topic}" || got.Selectors.SourceMatch != "^guide/" {
		t.Fatalf("parseGenerateReportArgs opts = %+v", got)
	}
	if len(gotPos) != 1 || gotPos[0] != "nb" {
		t.Fatalf("parseGenerateReportArgs positional = %q, want [nb]", gotPos)
	}
}

func TestParseCreateReportArgsUsesGlobalSelectors(t *testing.T) {
	inv, err := parseInvocation([]string{"--source-match", "^spec/", "create-report", "nb", "brief"}, nil, nil, os.Stderr)
	if err != nil {
		t.Fatalf("parseInvocation: %v", err)
	}
	got, gotPos, err := parseCreateReportArgsWithOptions(inv.args, inv.globals)
	if err != nil {
		t.Fatalf("parseCreateReportArgs: %v", err)
	}
	if got.Selectors.SourceMatch != "^spec/" {
		t.Fatalf("create-report selectors = %+v, want source-match from globals", got.Selectors)
	}
	if strings.Join(gotPos, ",") != "nb,brief" {
		t.Fatalf("create-report positional = %q, want [nb brief]", gotPos)
	}
}

func TestParseCreateReportArgsLocalSelectors(t *testing.T) {
	got, gotPos, err := parseCreateReportArgsWithOptions([]string{"nb", "brief", "--source-ids", "src-1,src-2", "desc"}, globalOptions{})
	if err != nil {
		t.Fatalf("parseCreateReportArgs: %v", err)
	}
	if got.Selectors.SourceIDs != "src-1,src-2" {
		t.Fatalf("create-report selectors = %+v, want source ids", got.Selectors)
	}
	if strings.Join(gotPos, ",") != "nb,brief,desc" {
		t.Fatalf("create-report positional = %q, want [nb brief desc]", gotPos)
	}
}

func TestParseChatShowArgsResolveCitationsCompatibility(t *testing.T) {
	got, gotPos, err := parseChatShowArgsWithOptions([]string{"nb", "conv", "--resolve-citations"}, globalOptions{})
	if err != nil {
		t.Fatalf("parseChatShowArgs: %v", err)
	}
	if !got.ResolveCitations {
		t.Fatalf("chat show resolve citations = false, want true")
	}
	if strings.Join(gotPos, ",") != "nb,conv" {
		t.Fatalf("chat show positional = %q, want [nb conv]", gotPos)
	}
}

func TestSaveChatSessionWritesConversationFile(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	session := &ChatSession{
		NotebookID:     "nb",
		ConversationID: "12345678-1234-1234-1234-123456789abc",
		Messages:       []ChatMessage{{Role: "user", Content: "hello"}},
	}
	if err := saveChatSession(session); err != nil {
		t.Fatalf("saveChatSession: %v", err)
	}
	if _, err := os.Stat(getChatSessionPath("nb")); err != nil {
		t.Fatalf("legacy session file missing: %v", err)
	}
	if _, err := os.Stat(getChatSessionPathForConv("nb", session.ConversationID)); err != nil {
		t.Fatalf("conversation session file missing: %v", err)
	}
}

func TestLoadChatSessionByConversationPrefixFallsBackToLegacy(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	session := `{
  "notebook_id": "nb",
  "conversation_id": "abcdef12-3456-7890-abcd-ef1234567890",
  "messages": [{"role": "assistant", "content": "smoke ok"}]
}`
	path := filepath.Join(home, ".nlm")
	if err := os.MkdirAll(path, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(path, "chat-nb.json"), []byte(session), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := loadChatSessionByConversation("nb", "abcdef12")
	if err != nil {
		t.Fatalf("loadChatSessionByConversation: %v", err)
	}
	if got.ConversationID != "abcdef12-3456-7890-abcd-ef1234567890" {
		t.Fatalf("conversation = %q", got.ConversationID)
	}
	if len(got.Messages) != 1 || !strings.Contains(got.Messages[0].Content, "smoke") {
		t.Fatalf("messages = %+v", got.Messages)
	}
}
