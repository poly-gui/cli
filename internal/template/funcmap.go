package template

import (
	"github.com/iancoleman/strcase"
	gotemplate "text/template"
)

var FuncMap = gotemplate.FuncMap{
	"ToKebab": strcase.ToKebab,
}
