package frame

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/laher/mergefs"
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

func (fa *FrameApplication) Templates(templateFS fs.FS, manifest TemplateCollection) *FrameApplication {
	var (
		err            error
		tmplDefinition Template
	)

	// fa.templateFS = templateFS
	fa.templateFS = mergefs.Merge(templateFS, internalTemplatesFS)
	fa.templates = make(map[string]*template.Template)

	manifest = fa.registerInternalTemplates(manifest)

	for _, tmplDefinition = range manifest {
		var parsedTemplate *template.Template

		tmplPath := filepath.Join("templates", tmplDefinition.Name)
		layoutPath := filepath.Join("templates", tmplDefinition.UseLayout)

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

func (fa *FrameApplication) registerInternalTemplates(manifest TemplateCollection) TemplateCollection {
	manifest = append(manifest, Template{Name: "account-pending.tmpl", IsLayout: false, UseLayout: "layout.tmpl"})
	manifest = append(manifest, Template{Name: "login.tmpl", IsLayout: false, UseLayout: "layout.tmpl"})
	manifest = append(manifest, Template{Name: "unexpected-error.tmpl", IsLayout: false, UseLayout: "layout.tmpl"})

	return manifest
}

func (fa *FrameApplication) setupInternalTemplate(templateString, layoutName, newTemplateName string) {
	var (
		err error
	)

	fa.Lock()
	defer fa.Unlock()

	// Validate the layout template exists first
	layoutFileName := fmt.Sprintf("%s.tmpl", layoutName)

	if _, ok := fa.templates[layoutFileName]; !ok {
		fa.Logger.Fatalf("The layout template name '%s' does not exist", layoutName)
	}

	// Render the login template and add it to the template cache
	// newTemplate, err := template.New(newTemplateName).Parse(templateString)
	// newTemplate, err := fa.templates[layoutFileName].Parse(templateString)
	newTemplate, err := fa.templates[layoutFileName].New(newTemplateName).Parse(templateString)

	if err != nil {
		fa.Logger.WithError(err).Fatalf("error parsing %s template", newTemplateName)
	}

	fa.templates[newTemplateName] = newTemplate
	fmt.Printf("\n%+v\n", fa.templates)
}
