package cli

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"poly-cli/internal/poly"
	"poly-cli/internal/template"
	"runtime"
	"sync"
)

// Generate is run when `poly generate` is invoked.
func Generate() error {
	os.Args = append(os.Args[:1], os.Args[2:]...)

	project, err := parseFromArgs()
	if err != nil {
		return err
	}

	err = os.MkdirAll(project.FullPath, os.ModePerm)
	if err != nil {
		return err
	}

	git, err := exec.LookPath("git")
	if err != nil {
		return err
	}
	cmd := exec.Command(git, "init")
	cmd.Dir = project.FullPath
	err = cmd.Run()
	if err != nil {
		return err
	}

	errs := make([]error, 3)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		errs[0] = generateMacOSSource(*project)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		errs[1] = generateGTKSource(*project)
		wg.Done()
	}()

	wg.Wait()

	err = generatePortableLayerSource(*project)

	return errors.Join(errs...)
}

func generateMacOSSource(project poly.ProjectDescription) error {
	log.Println("Generating macOS project...")

	o, err := filepath.Abs(filepath.Join(project.FullPath, defaultMacOSFolderName))
	if err != nil {
		return err
	}

	err = template.GenerateTemplates(template.AppKitSourceFiles, o, project)
	if err != nil {
		return err
	}

	if runtime.GOOS == "darwin" {
		xcodegen, err := exec.LookPath("xcodegen")
		if err != nil {
			return err
		}
		cmd := exec.Command(xcodegen, "generate")
		cmd.Dir = o

		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func generateGTKSource(project poly.ProjectDescription) error {
	log.Println("Generating GTK project...")

	git, err := exec.LookPath("git")
	if err != nil {
		return err
	}

	o, err := filepath.Abs(filepath.Join(project.FullPath, defaultGtkFolderName))
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Join(o), os.ModePerm)
	if err != nil {
		return err
	}

	libp := filepath.Join(project.FullPath, "gtk", "lib")

	err = os.MkdirAll(libp, os.ModePerm)
	if err != nil {
		return err
	}

	cmd := exec.Command(git, "submodule", "add", gtkpolyGitURL)
	cmd.Dir = libp
	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command(git, "submodule", "update", "--init", "--recursive")
	cmd.Dir = project.FullPath
	err = cmd.Run()
	if err != nil {
		return err
	}

	return template.GenerateTemplates(template.GTKSourceFiles, o, project)
}

func parseFromArgs() (*poly.ProjectDescription, error) {
	var outputPath string
	var projectName string
	var packageName string
	var debugWorkspacePath string
	var language string

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	flag.StringVar(&debugWorkspacePath, "debug-workspace", "", "")
	flag.StringVar(&outputPath, "output", cwd, "Where the project should be created in. Defaults to the current working directory.")
	flag.StringVar(&projectName, "name", defaultProjectName, "The name for the application. Default is "+defaultProjectName+".")
	flag.StringVar(&packageName, "package", defaultPackageName, "The package name/bundle ID for the application. Default is "+defaultPackageName+".")
	flag.StringVar(&language, "language", defaultLanguage, "The programming language the source code of the portable layer should be using.")

	flag.Parse()

	if !poly.IsLanguageSupported(language) {
		return nil, errors.New(fmt.Sprintf("Unsupported language: %v", language))
	}

	o, err := filepath.Abs(filepath.Join(outputPath, projectName))
	if err != nil {
		return nil, err
	}

	projectDescription := poly.ProjectDescription{
		DebugWorkspacePath: debugWorkspacePath,
		FullPath:           o,
		AppName:            projectName,
		OrganizationName:   "",
		PackageName:        packageName,
		Language:           poly.SupportedLanguage(language),
	}

	return &projectDescription, nil
}

func generatePortableLayerSource(project poly.ProjectDescription) error {
	switch project.Language {
	case poly.SupportedLanguageTypeScript:
		return generateTypeScriptSource(project)
	default:
		return errors.New(fmt.Sprintf("Unsupported language %v\n", project.Language))
	}
}
