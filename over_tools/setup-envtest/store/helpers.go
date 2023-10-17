// SPDX-License-Identifier: Apache-2.0
// Copyright 2021 The Kubernetes Authors

package store

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

// DefaultStoreDir returns the default location for the store.
// It's dependent on operating system:
//
// - Windows: %LocalAppData%\kubebuilder-over_envtest
// - OSX: ~/Library/Application Support/io.kubebuilder.over_envtest
// - Others: ${XDG_DATA_HOME:-~/.local/share}/kubebuilder-over_envtest
//
// Otherwise, it errors out.  Note that these paths must not be relied upon
// manually.
func DefaultStoreDir() (string, error) {
	var baseDir string

	// find the base data directory
	switch runtime.GOOS {
	case "windows":
		baseDir = os.Getenv("LocalAppData")
		if baseDir == "" {
			return "", errors.New("%LocalAppData% is not defined")
		}
	case "darwin", "ios":
		homeDir := os.Getenv("HOME")
		if homeDir == "" {
			return "", errors.New("$HOME is not defined")
		}
		baseDir = filepath.Join(homeDir, "Library/Application Support")
	default:
		baseDir = os.Getenv("XDG_DATA_HOME")
		if baseDir == "" {
			homeDir := os.Getenv("HOME")
			if homeDir == "" {
				return "", errors.New("neither $XDG_DATA_HOME nor $HOME are defined")
			}
			baseDir = filepath.Join(homeDir, ".local/share")
		}
	}

	// append our program-specific dir to it (OSX has a slightly different
	// convention so try to follow that).
	switch runtime.GOOS {
	case "darwin", "ios":
		return filepath.Join(baseDir, "io.kubebuilder.over_envtest"), nil
	default:
		return filepath.Join(baseDir, "kubebuilder-over_envtest"), nil
	}
}
