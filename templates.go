package frame

import (
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
)

type Template struct {
	Name      string
	IsLayout  bool
	UseLayout string
}

type TemplateCollection []Template

func (fa *FrameApplication) RenderTemplate(w http.ResponseWriter, name string, data map[string]interface{}) {
	var (
		err  error
		tmpl *template.Template
		ok   bool
	)

	if tmpl, ok = fa.templates[name]; !ok {
		fa.Logger.Fatalf("template '%s' not found!", name)
	}

	if err = tmpl.Execute(w, data); err != nil {
		fa.Logger.WithError(err).Fatalf("error rendering '%s'", name)
	}
}

func (fa *FrameApplication) Templates(templateFS fs.FS, rootPath string, manifest TemplateCollection) *FrameApplication {
	var (
		err            error
		tmplDefinition Template
	)

	fa.templateFS = templateFS
	fa.templates = make(map[string]*template.Template)

	for _, tmplDefinition = range manifest {
		var parsedTemplate *template.Template

		tmplPath := filepath.Join(rootPath, tmplDefinition.Name)
		layoutPath := filepath.Join(rootPath, tmplDefinition.UseLayout)

		if tmplDefinition.IsLayout {
			if parsedTemplate, err = template.ParseFS(fa.templateFS, tmplPath); err != nil {
				fa.Logger.WithError(err).Fatalf("error parsing layout '%s'. shutting down", tmplDefinition.Name)
			}
		} else {
			if parsedTemplate, err = template.ParseFS(fa.templateFS, tmplPath, layoutPath); err != nil {
				fa.Logger.WithError(err).Fatalf("error parsing template '%s' with layout '%s'. shutting down", tmplDefinition.Name, tmplDefinition.UseLayout)
			}
		}

		fa.templates[tmplDefinition.Name] = parsedTemplate
	}

	return fa
}
