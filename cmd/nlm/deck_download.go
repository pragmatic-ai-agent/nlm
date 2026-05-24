package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type deckDownloadOptions struct {
	ArtifactID string
	Format     string
	Output     string
}

func printDeckDownloadUsage(cmdName string) {
	fmt.Fprintf(os.Stderr, "Usage: nlm %s <notebook-id> --id <artifact-id> [--format pdf|pptx] [--output file]\n\n", cmdName)
	fmt.Fprintln(os.Stderr, "The current CLI cannot fetch slide deck files directly. It prints the")
	fmt.Fprintln(os.Stderr, "NotebookLM browser URL so you can use the web UI download menu.")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Flags:")
	fmt.Fprintln(os.Stderr, "  --id, --artifact-id <id>  Slide deck artifact ID")
	fmt.Fprintln(os.Stderr, "  --format, -f <format>    Desired browser export format: pdf or pptx")
	fmt.Fprintln(os.Stderr, "  --output, -o <file>      Intended output filename for your browser save dialog")
}

func validateDeckDownloadArgs(cmdName string, args []string) error {
	return validateDeckDownloadArgsWithOptions(cmdName, args, globalOptions{})
}

func validateDeckDownloadArgsWithOptions(cmdName string, args []string, _ globalOptions) error {
	if _, _, err := parseDeckDownloadArgs(args); err != nil {
		fmt.Fprintf(os.Stderr, "nlm: %s: %v\n\n", cmdName, err)
		printDeckDownloadUsage(cmdName)
		return errBadArgs
	}
	return nil
}

func parseDeckDownloadArgs(args []string) (deckDownloadOptions, string, error) {
	opts := deckDownloadOptions{Format: "pdf"}
	flags := flag.NewFlagSet("deck-download", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&opts.ArtifactID, "id", "", "")
	flags.StringVar(&opts.ArtifactID, "artifact-id", "", "")
	flags.StringVar(&opts.Format, "format", opts.Format, "")
	flags.StringVar(&opts.Format, "f", opts.Format, "")
	flags.StringVar(&opts.Output, "output", "", "")
	flags.StringVar(&opts.Output, "o", "", "")

	flagArgs, positional, err := splitCommandFlags(args, map[string]bool{
		"id":          true,
		"artifact-id": true,
		"format":      true,
		"f":           true,
		"output":      true,
		"o":           true,
	}, nil)
	if err != nil {
		return opts, "", err
	}
	if err := flags.Parse(flagArgs); err != nil {
		return opts, "", err
	}
	if len(positional) != 1 {
		return opts, "", fmt.Errorf("requires exactly one notebook id")
	}
	if opts.ArtifactID == "" {
		return opts, "", fmt.Errorf("missing --id <artifact-id>")
	}
	switch opts.Format {
	case "pdf", "pptx":
	default:
		return opts, "", fmt.Errorf("unsupported format %q (want pdf or pptx)", opts.Format)
	}
	return opts, positional[0], nil
}

func runDeckDownloadFallback(args []string) error {
	return runDeckDownloadFallbackWithOptions(args, globalOptions{})
}

func runDeckDownloadFallbackWithOptions(args []string, _ globalOptions) error {
	opts, notebookID, err := parseDeckDownloadArgs(args)
	if err != nil {
		return err
	}
	u := notebookBrowserURL(notebookID)
	fmt.Println(u)
	fmt.Fprintf(os.Stderr, "Direct slide deck download is not implemented for artifact %s.\n", opts.ArtifactID)
	fmt.Fprintf(os.Stderr, "Open %s in a browser and use NotebookLM's download menu for %s export.\n", u, opts.Format)
	if opts.Output != "" {
		fmt.Fprintf(os.Stderr, "Requested output filename: %s\n", opts.Output)
	}
	return fmt.Errorf("direct slide deck download unavailable; open %s in a browser", u)
}
