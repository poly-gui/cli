package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"poly-cli/internal/poly"
	"poly-cli/internal/template"
	"strconv"
	"strings"
)

func generateTypeScriptSource(project poly.ProjectDescription) error {
	pms, err := findPackageManagers()
	if err != nil {
		return err
	}
	if len(pms) == 0 {
		return errors.New("no package manager was found! make sure at least one package manager (e.g. npm) is installed on your computer and is accessible in PATH.")
	}

	var pm string
	if len(pms) == 1 {
		for _, path := range pms {
			pm = path
		}
	} else {
		reader := bufio.NewReader(os.Stdin)

		fmt.Println("Multiple package manager found:")

		pathList := make([]string, 0, len(pms))
		i := 0
		for name, path := range pms {
			pathList = append(pathList, path)
			fmt.Printf("(%d) %v\n", i, name)
			i++
		}

		for {
			fmt.Print("Which one would you like to use? Enter the number: ")
			s, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			i, err := strconv.Atoi(strings.TrimSpace(s))
			if err != nil || i >= len(pathList) {
				continue
			}
			pm = pathList[i]
			break
		}
	}

	o, err := filepath.Abs(filepath.Join(project.FullPath, defaultAppSrcFolderName))
	if err != nil {
		return err
	}

	err = template.GenerateTemplates(template.TSSourceFiles, o, project)
	if err != nil {
		return err
	}

	cmd := exec.Command(pm, "install")
	cmd.Dir = o
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func findPackageManagers() (map[string]string, error) {
	pms := map[string]string{}

	npm, err := exec.LookPath("npm")
	if err == nil {
		pms["npm"] = npm
	}

	pnpm, err := exec.LookPath("pnpm")
	if err == nil {
		pms["pnpm"] = pnpm
	}

	yarn, err := exec.LookPath("yarn")
	if err == nil {
		pms["yarn"] = yarn
	}

	bun, err := exec.LookPath("bun")
	if err == nil {
		pms["bun"] = bun
	}

	return pms, nil
}
