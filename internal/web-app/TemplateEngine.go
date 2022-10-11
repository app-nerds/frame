package webapp

import (
	"html/template"
	"net/http"
	"path/filepath"

	pkgwebapp "github.com/app-nerds/frame/pkg/web-app"
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
		tmplDefinition pkgwebapp.Template
	)

	wa.templateFS = mergefs.Merge(wa.templateFS, wa.internalTemplateFS)
	wa.templateManifest = wa.registerInternalTemplates()

	for _, tmplDefinition = range wa.templateManifest {
		var parsedTemplate *template.Template

		tmplPath := filepath.Join("templates", tmplDefinition.Name)
		layoutPath := filepath.Join("templates", tmplDefinition.UseLayout)

		if tmplDefinition.IsLayout {
			if parsedTemplate, err = template.ParseFS(wa.templateFS, tmplPath); err != nil {
				wa.logger.WithError(err).Fatalf("error parsing layout '%s'. shutting down", tmplDefinition.Name)
			}
		} else {
			if parsedTemplate, err = template.ParseFS(wa.templateFS, tmplPath, layoutPath); err != nil {
				wa.logger.WithError(err).Fatalf("error parsing template '%s' with layout '%s'. shutting down", tmplDefinition.Name, tmplDefinition.UseLayout)
			}
		}

		wa.templates[tmplDefinition.Name] = parsedTemplate
	}

	/*
	 * Log out some debug information in development mode
	 */
	if wa.debug {
		for k, v := range wa.templates {
			wa.logger.WithFields(logrus.Fields{
				"templateName": k,
				"value":        v,
			}).Info("template captured")
		}
	}
}

func (wa *WebApp) registerAdminTemplates() pkgwebapp.TemplateCollection {
	manifest := pkgwebapp.TemplateCollection{}
	manifest = append(manifest, pkgwebapp.Template{Name: "admin-layout.tmpl", IsLayout: true})
	manifest = append(manifest, pkgwebapp.Template{Name: "admin-dashboard.tmpl", IsLayout: false, UseLayout: "admin-layout.tmpl"})
	manifest = append(manifest, pkgwebapp.Template{Name: "admin-members-manage.tmpl", IsLayout: false, UseLayout: "admin-layout.tmpl"})

	return manifest
}

func (wa *WebApp) registerInternalTemplates() pkgwebapp.TemplateCollection {
	wa.templateManifest = append(wa.templateManifest, pkgwebapp.Template{Name: "account-pending.tmpl", IsLayout: false, UseLayout: "layout.tmpl"})
	wa.templateManifest = append(wa.templateManifest, pkgwebapp.Template{Name: "login.tmpl", IsLayout: false, UseLayout: "layout.tmpl"})
	wa.templateManifest = append(wa.templateManifest, pkgwebapp.Template{Name: "unexpected-error.tmpl", IsLayout: false, UseLayout: "layout.tmpl"})
	wa.templateManifest = append(wa.templateManifest, pkgwebapp.Template{Name: "sign-up.tmpl", IsLayout: false, UseLayout: "layout.tmpl"})

	return wa.templateManifest
}

func (wa *WebApp) setupAdminTemplates() {
	var (
		err            error
		tmplDefinition pkgwebapp.Template
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
