package frame

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/laher/mergefs"
	"github.com/sirupsen/logrus"
)

func (wa *WebApp) RenderTemplate(w http.ResponseWriter, name string, data interface{}) {
	var (
		err  error
		tmpl *template.Template
		ok   bool
	)

	if tmpl, ok = wa.templates[name]; !ok {
		wa.logger.Fatalf("template '%s' not found!", name)
	}

	if err = tmpl.Execute(w, data); err != nil {
		wa.logger.WithError(err).Fatalf("error rendering '%s'", name)
	}
}

func (wa *WebApp) setupTemplateEngine() {
	var (
		err            error
		tmplDefinition Template
	)

	wa.templateFS = mergefs.Merge(wa.templateFS, wa.internalTemplateFS)
	wa.templateManifest = wa.registerInternalTemplates()

	templateFuncs := template.FuncMap{
		"IsSet": wa.templateFuncIsSet,
	}

	for _, tmplDefinition = range wa.templateManifest {
		var parsedTemplate *template.Template

		tmplPath := filepath.Join("frontend-templates", tmplDefinition.Name)
		layoutPath := filepath.Join("frontend-templates", tmplDefinition.UseLayout)

		if tmplDefinition.IsLayout {
			if parsedTemplate, err = template.New(tmplDefinition.Name).Funcs(templateFuncs).ParseFS(wa.templateFS, tmplPath); err != nil {
				wa.logger.WithError(err).Fatalf("error parsing layout '%s'. shutting down", tmplDefinition.Name)
			}
		} else {
			if tmplDefinition.UseLayout == "" {
				wa.logger.Fatalf("The template '%s' must have a layout specified.", tmplDefinition.Name)
			}

			if parsedTemplate, err = template.New(tmplDefinition.Name).Funcs(templateFuncs).ParseFS(wa.templateFS, tmplPath, layoutPath); err != nil {
				wa.logger.WithError(err).Fatalf("error parsing template '%s' with layout '%s'. shutting down", tmplDefinition.Name, tmplDefinition.UseLayout)
			}
		}

		wa.templates[tmplDefinition.Name] = parsedTemplate
	}

	/*
	 * Log out some debug information in development mode
	 */
	if wa.debug {
		for k := range wa.templates {
			wa.logger.WithFields(logrus.Fields{
				"templateName": k,
			}).Info("template captured")
		}
	}
}

func (wa *WebApp) setupAdminTemplates() {
	var (
		err            error
		tmplDefinition Template
	)

	manifest := wa.registerAdminTemplates()

	for _, tmplDefinition = range manifest {
		var parsedTemplate *template.Template

		tmplPath := filepath.Join("admin-templates", tmplDefinition.Name)
		layoutPath := filepath.Join("admin-templates", tmplDefinition.UseLayout)

		if tmplDefinition.IsLayout {
			if parsedTemplate, err = template.ParseFS(wa.adminTemplateFS, tmplPath); err != nil {
				wa.logger.WithError(err).Fatalf("error parsing admin layout '%s'. shutting down", tmplDefinition.Name)
			}
		} else {
			if parsedTemplate, err = template.ParseFS(wa.adminTemplateFS, tmplPath, layoutPath); err != nil {
				wa.logger.WithError(err).Fatalf("error parsing admin template '%s' with layout '%s'. shutting down", tmplDefinition.Name, tmplDefinition.UseLayout)
			}
		}

		wa.templates[tmplDefinition.Name] = parsedTemplate
	}
}
