package frame

import (
	"html/template"
	"io/fs"
	"path/filepath"
)

type Template struct {
	Name      string
	IsLayout  bool
	UseLayout string
}

type TemplateCollection []Template

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
				fa.logger.WithError(err).Fatalf("error parsing layout '%s'. shutting down", tmplDefinition.Name)
			}
		} else {
			if parsedTemplate, err = template.ParseFS(fa.templateFS, tmplPath, layoutPath); err != nil {
				fa.logger.WithError(err).Fatalf("error parsing template '%s' with layout '%s'. shutting down", tmplDefinition.Name, tmplDefinition.UseLayout)
			}
		}

		fa.templates[tmplDefinition.Name] = parsedTemplate
	}

	return fa
}
