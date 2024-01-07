package cli

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"poly-cli/internal/poly"
	"poly-cli/internal/template"
	gotemplate "text/template"
)

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

	err = generateMacOSSource(projectDescription)
	if err != nil {
		return err
	}

	err = generatePortableLayerSource(projectDescription)
	if err != nil {
		return err
	}

	return nil
}

func generateMacOSSource(project poly.ProjectDescription) error {
	log.Println("Generating macOS project...")

	xcodegen, err := exec.LookPath("xcodegen")
	if err != nil {
		return err
	}

	o, err := filepath.Abs(filepath.Join(project.FullPath, defaultMacOSFolderName))
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Join(o, project.AppName), os.ModePerm)
	if err != nil {
		return err
	}

	tmpl, err := gotemplate.New("XcodeGenSpec").Parse(template.XcodeGenSpec)
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(o, xcodeProjectSpecName))
	if err != nil {
		return err
	}

	err = tmpl.Execute(f, project)
	if err != nil {
		return err
	}

	_ = f.Close()

	f, err = os.Create(filepath.Join(o, project.AppName, "AppDelegate.swift"))
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(template.AppDelegate))
	if err != nil {
		return err
	}
	_ = f.Close()

	f, err = os.Create(filepath.Join(o, project.AppName, "main.swift"))
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(template.SwiftMainFile))
	if err != nil {
		return err
	}
	_ = f.Close()

	cmd := exec.Command(xcodegen, "generate")
	cmd.Dir = o

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func generatePortableLayerSource(project poly.ProjectDescription) error {
	bun, err := exec.LookPath("bun")
	if err != nil {
		return err
	}

	o, err := filepath.Abs(filepath.Join(project.FullPath, defaultAppSrcFolderName))
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Join(o, "src"), os.ModePerm)
	if err != nil {
		return err
	}

	tmpl, err := gotemplate.New("PackageJSON").Funcs(template.FuncMap).Parse(template.PackageJSON)
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(o, "package.json"))
	if err != nil {
		return err
	}
	err = tmpl.Execute(f, project)
	if err != nil {
		return err
	}
	_ = f.Close()

	f, err = os.Create(filepath.Join(o, "tsconfig.json"))
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(template.TSConfigJSON))
	if err != nil {
		return err
	}
	_ = f.Close()

	f, err = os.Create(filepath.Join(o, ".gitignore"))
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(template.TSGitIgnore))
	if err != nil {
		return err
	}
	_ = f.Close()

	f, err = os.Create(filepath.Join(o, "src", "main.ts"))
	if err != nil {
		return err
	}
	tmpl, err = gotemplate.New("TSMainFile").Parse(template.TSMainFile)
	if err != nil {
		return err
	}
	err = tmpl.Execute(f, project)
	if err != nil {
		return err
	}
	_ = f.Close()

	cmd := exec.Command(bun, "install")
	cmd.Dir = o
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
