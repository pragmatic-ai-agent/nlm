package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/tmc/nlm/internal/notebooklm/api"
)

type updateArtifactOptions struct {
	Name string
}

func validateUpdateArtifactArgsWithOptions(cmdName string, args []string, globals globalOptions) error {
	if _, _, err := parseUpdateArtifactArgsWithOptions(args, globals); err != nil {
		return fmt.Errorf("%w: %v", errBadArgs, err)
	}
	return nil
}

func parseUpdateArtifactArgsWithOptions(args []string, globals globalOptions) (updateArtifactOptions, string, error) {
	opts := updateArtifactOptions{Name: globals.sourceName}
	flags := flag.NewFlagSet("update-artifact", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&opts.Name, "name", opts.Name, "")
	flags.StringVar(&opts.Name, "n", opts.Name, "")

	flagArgs, positional, err := splitCommandFlags(args, map[string]bool{
		"name": true, "n": true,
	}, nil)
	if err != nil {
		return opts, "", err
	}
	if err := flags.Parse(flagArgs); err != nil {
		return opts, "", err
	}
	if len(positional) < 1 || len(positional) > 2 {
		return opts, "", fmt.Errorf("want artifact id and optional title")
	}
	if len(positional) == 2 {
		opts.Name = positional[1]
	}
	if opts.Name == "" {
		return opts, "", fmt.Errorf("provide new title as second arg or --name flag")
	}
	return opts, positional[0], nil
}

func runUpdateArtifactWithOptions(c *api.Client, args []string, globals globalOptions) error {
	opts, artifactID, err := parseUpdateArtifactArgsWithOptions(args, globals)
	if err != nil {
		return err
	}
	return renameArtifact(c, artifactID, opts.Name)
}
