package main

import "testing"

func TestParseUpdateArtifactArgsWithOptions(t *testing.T) {
	opts, artifactID, err := parseUpdateArtifactArgsWithOptions([]string{"art-1", "--name", "New"}, globalOptions{})
	if err != nil {
		t.Fatalf("parseUpdateArtifactArgsWithOptions: %v", err)
	}
	if artifactID != "art-1" || opts.Name != "New" {
		t.Fatalf("artifactID, opts = %q, %+v; want art-1 New", artifactID, opts)
	}

	opts, artifactID, err = parseUpdateArtifactArgsWithOptions([]string{"art-2"}, globalOptions{sourceName: "FromGlobal"})
	if err != nil {
		t.Fatalf("parseUpdateArtifactArgsWithOptions inherited: %v", err)
	}
	if artifactID != "art-2" || opts.Name != "FromGlobal" {
		t.Fatalf("artifactID, opts = %q, %+v; want art-2 FromGlobal", artifactID, opts)
	}
}
