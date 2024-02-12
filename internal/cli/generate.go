package cli

import (
	"errors"
	"flag"
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

	var outputPath string
	var projectName string
	var packageName string
	var debugMode bool

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	flag.BoolVar(&debugMode, "debug", false, "Enables debug mode.")
	flag.StringVar(&outputPath, "output", cwd, "Where the project should be created in. Defaults to the current working directory.")
	flag.StringVar(&projectName, "name", defaultProjectName, "The name for the application. Default is "+defaultProjectName+".")
	flag.StringVar(&packageName, "package", defaultPackageName, "The package name/bundle ID for the application. Default is "+defaultPackageName+".")

	flag.Parse()

	o, err := filepath.Abs(filepath.Join(outputPath, projectName))
	if err != nil {
		return err
	}

	projectDescription := poly.ProjectDescription{
		DebugMode:        debugMode,
		FullPath:         o,
		AppName:          projectName,
		OrganizationName: "",
		PackageName:      packageName,
	}

	err = os.MkdirAll(o, os.ModePerm)
	if err != nil {
		return err
	}

	git, err := exec.LookPath("git")
	if err != nil {
		return err
	}
	cmd := exec.Command(git, "init")
	cmd.Dir = o
	err = cmd.Run()
	if err != nil {
		return err
	}

	errs := make([]error, 3)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		errs[0] = generateMacOSSource(projectDescription)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		errs[1] = generateGTKSource(projectDescription)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		errs[2] = generatePortableLayerSource(projectDescription)
		wg.Done()
	}()

	wg.Wait()

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

func generatePortableLayerSource(project poly.ProjectDescription) error {
	pnpm, err := exec.LookPath("pnpm")
	if err != nil {
		return err
	}

	o, err := filepath.Abs(filepath.Join(project.FullPath, defaultAppSrcFolderName))
	if err != nil {
		return err
	}

	err = template.GenerateTemplates(template.TSSourceFiles, o, project)
	if err != nil {
		return err
	}

	cmd := exec.Command(pnpm, "install")
	cmd.Dir = o
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
