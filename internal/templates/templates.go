package templates

import (
	"html/template"
	"io"
	"path"
)

type Templates struct {
	templateDir string
	baseLayout  string
	templates   []string
}

func New(templateDir string, baseLayout string) *Templates {
	return &Templates{
		templateDir: templateDir,
		baseLayout:  baseLayout,
	}
}

func (t *Templates) Add(template string) *Templates {
	t.templates = append(t.templates, path.Join(t.templateDir, template))
	return t
}

func (t *Templates) Parse() func(io.Writer, interface{}) error {
	parsedTemplates := template.Must(template.ParseFiles(t.templates...))
	return func(writer io.Writer, data interface{}) error {
		return parsedTemplates.ExecuteTemplate(writer, t.baseLayout, data)
	}
}
