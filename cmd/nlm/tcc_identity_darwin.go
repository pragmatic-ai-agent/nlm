//go:build darwin

package main

import (
	"fmt"
	"os"

	"github.com/tmc/macgo"
	"golang.org/x/term"
)

var startTCCIdentityApp = func(cfg *macgo.Config) error {
	return macgo.Start(cfg)
}

var tccIdentityTTY = term.IsTerminal

func prepareTCCIdentity(appName, bundleID, usage string, perms ...macgo.Permission) error {
	if tccIdentityTTY(int(os.Stdin.Fd())) {
		_ = os.Setenv("MACGO_TTY_PASSTHROUGH", "1")
	}

	cfg := macgo.NewConfig().
		WithAppName(appName).
		WithBundleID(bundleID).
		WithPermissions(perms...).
		WithDevMode().
		WithAdHocSign().
		WithUIMode(macgo.UIModeBackground)
	if usage != "" {
		cfg.WithMicrophoneUsage(usage)
	}
	if debug {
		cfg.WithDebug()
	}
	if err := startTCCIdentityApp(cfg); err != nil {
		return fmt.Errorf("prepare TCC identity: %w", err)
	}
	return nil
}
