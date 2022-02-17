//go:build mage

package main

import (
	"os"
	"runtime"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Clear deletes the app binary
func Clear() error {
	return sh.Run("rm", "app", "-f")
}

// Build compiles the app
func Build() error {
	mg.Deps(Clear)

	if err := sh.Run("go", "mod", "download"); err != nil {
		return err
	}

	env := map[string]string{
		"GOOS":   runtime.GOOS,
		"GOARCH": runtime.GOARCH,
	}
	_, err := sh.Exec(env, os.Stdout, os.Stderr, "go", "build", "-ldflags="+"-w -s", "-o", "app", "./cmd/server")

	return err
}

// Test runs all tests inside /pkg
func Test() error {
	_, err := sh.Exec(nil, os.Stdout, os.Stderr, "go", "test", "./pkg/...", "-v", "-cover")
	return err
}

// Docker creates Docker image for the app
func Docker() error {
	mg.Deps(Clear)

	_, err := sh.Exec(nil, os.Stdout, os.Stderr, "docker", "build", "-t", "fpl-live-tracker", ".")
	sh.Run("docker", "image", "prune", "-f")
	return err
}
