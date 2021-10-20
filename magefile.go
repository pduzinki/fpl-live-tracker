//go:build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

func Clear() error {
	return sh.Run("rm", "app", "-f")
}

func Build() error {
	mg.Deps(Clear)

	if err := sh.Run("go", "mod", "download"); err != nil {
		return err
	}

	env := map[string]string{
		"GOOS":   "linux",
		"GOARCH": "amd64",
	}
	return sh.RunWith(env, "go", "build", "-ldflags="+"-w -s", "-o", "app", "./cmd/server")
}

func Docker() error {
	mg.Deps(Clear)

	return sh.Run("docker", "build", "-t", "fpl-live-tracker", ".")
}
