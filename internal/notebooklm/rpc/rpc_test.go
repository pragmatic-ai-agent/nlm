package rpc

import "testing"

func TestNewWithConfigUsesAuthUserFromEnv(t *testing.T) {
	t.Setenv("NLM_AUTHUSER", "2")

	client := NewWithConfig("token", "cookies", ServiceConfig{
		URLParams: map[string]string{
			"authuser": "1",
		},
	})

	if got := client.Config.URLParams["authuser"]; got != "2" {
		t.Fatalf("authuser URL param = %q, want 2", got)
	}
	if got := client.Config.Headers["x-goog-authuser"]; got != "2" {
		t.Fatalf("x-goog-authuser header = %q, want 2", got)
	}
}
