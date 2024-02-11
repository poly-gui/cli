package template

import (
	"errors"
	"github.com/iancoleman/strcase"
	"os"
	"path/filepath"
	"poly-cli/internal/poly"
	"strings"
	"sync"
	gotemplate "text/template"
)

type templateFile struct {
	// FilePathRel is the path to the output file relative to the native source code directory such as <project>/macOS or <project>/gtk.
	// Use _APP_NAME_ as a placeholder for the app name if needed.
	FilePathRel string

	// Template is the go template of this file. Available functions:
	//   - ToKebab: converts the given string to kebab-case.
	Template string

	// TemplateName is an arbitrary name for the template for use with `templates.New("<the-name>")`
	TemplateName string
}

var funcMap = gotemplate.FuncMap{
	"ToKebab": strcase.ToKebab,
}

// GenerateTemplates generate the given list of template files inside the given directory for the given project.
// The directory must be specified in absolute path.
func GenerateTemplates(templates []templateFile, dir string, project poly.ProjectDescription) error {
	errs := make([]error, len(templates))
	var wg sync.WaitGroup

	for i, tmpl := range templates {
		tmpl, i := tmpl, i
		wg.Add(1)
		go func() {
			defer wg.Done()

			p := filepath.Join(dir, tmpl.FilePathRel)
			p = strings.Replace(p, "_APP_NAME_", project.AppName, -1)

			err := os.MkdirAll(filepath.Dir(p), os.ModePerm)
			if err != nil {
				errs[i] = err
				return
			}

			f, err := os.Create(p)
			if err != nil {
				errs[i] = err
				return
			}
			defer f.Close()

			t, err := gotemplate.New(tmpl.Template).Funcs(funcMap).Parse(tmpl.Template)
			if err != nil {
				errs[i] = err
				return
			}

			errs[i] = t.Execute(f, project)
		}()
	}

	wg.Wait()

	return errors.Join(errs...)
}
