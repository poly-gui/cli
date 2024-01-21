package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func Open(args ...string) error {
	switch args[0] {
	case "macos":
		return openMacOSProject()
	default:
		return errors.New(fmt.Sprintf("cannot open %v", args[0]))
	}
}

func openMacOSProject() error {
	if !strings.HasPrefix(runtime.GOOS, "darwin") {
		return errors.New("cannot open Xcode project on a non-macOS system")
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	if _, err := os.Stat(filepath.Join(wd, "macOS")); os.IsNotExist(err) {
		return errors.New("no macOS project found in the current working directory. make sure you are at the root of a Poly project")
	}

	cmd := exec.Command("xed", "macOS")
	return cmd.Run()
}
