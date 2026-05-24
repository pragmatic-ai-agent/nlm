package main

import "testing"

func TestParseAuthFlagsInterleaved(t *testing.T) {
	opts, remaining, err := parseAuthFlagsWithOptions([]string{
		"Work", "--debug", "--authuser", "2", "--cdp-url", "ws://localhost:9222",
	}, globalOptions{chromeProfile: "Default"})
	if err != nil {
		t.Fatalf("parseAuthFlagsWithOptions: %v", err)
	}
	if opts.ProfileName != "Default" {
		t.Fatalf("ProfileName = %q, want inherited Default", opts.ProfileName)
	}
	if !opts.Debug {
		t.Fatalf("Debug = false, want true")
	}
	if opts.AuthUser != "2" {
		t.Fatalf("AuthUser = %q, want 2", opts.AuthUser)
	}
	if opts.RemoteCDPURL != "ws://localhost:9222" {
		t.Fatalf("RemoteCDPURL = %q, want ws://localhost:9222", opts.RemoteCDPURL)
	}
	if len(remaining) != 1 || remaining[0] != "Work" {
		t.Fatalf("remaining = %v, want [Work]", remaining)
	}
}
