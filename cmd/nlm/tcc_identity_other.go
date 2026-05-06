//go:build !darwin

package main

func prepareTCCIdentity(appName, bundleID, usage string, perms ...interface{}) error {
	return nil
}
