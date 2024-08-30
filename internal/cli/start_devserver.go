package cli

import (
	"os"
	"poly-cli/internal/devserver"
)

func StartDevServer() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	return devserver.Start(devserver.Options{
		ProjectPathAbs: cwd,
	})
}
