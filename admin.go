package frame

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func (fa *FrameApplication) setupAdminTemplates() {
	var (
		err            error
		tmplDefinition Template
	)

	manifest := fa.registerAdminTemplates()

	for _, tmplDefinition = range manifest {
		var parsedTemplate *template.Template

		tmplPath := filepath.Join("admin-templates", tmplDefinition.Name)
		layoutPath := filepath.Join("admin-templates", tmplDefinition.UseLayout)

		if tmplDefinition.IsLayout {
			if parsedTemplate, err = template.ParseFS(adminTemplatesFS, tmplPath); err != nil {
				fa.Logger.WithError(err).Fatalf("error parsing admin layout '%s'. shutting down", tmplDefinition.Name)
			}
		} else {
			if parsedTemplate, err = template.ParseFS(adminTemplatesFS, tmplPath, layoutPath); err != nil {
				fa.Logger.WithError(err).Fatalf("error parsing admin template '%s' with layout '%s'. shutting down", tmplDefinition.Name, tmplDefinition.UseLayout)
			}
		}

		fa.templates[tmplDefinition.Name] = parsedTemplate
	}
}

func (fa *FrameApplication) registerAdminTemplates() TemplateCollection {
	manifest := TemplateCollection{}
	manifest = append(manifest, Template{Name: "admin-layout.tmpl", IsLayout: true})
	manifest = append(manifest, Template{Name: "admin-dashboard.tmpl", IsLayout: false, UseLayout: "admin-layout.tmpl"})
	manifest = append(manifest, Template{Name: "admin-members-manage.tmpl", IsLayout: false, UseLayout: "admin-layout.tmpl"})

	return manifest
}

func (fa *FrameApplication) setupAdminRoutes() {
	fa.router.HandleFunc("/admin", fa.handleAdminDashboard)
	fa.router.HandleFunc("/admin/members/manage", fa.handleAdminMembersManage).Methods(http.MethodGet)
	fa.router.HandleFunc("/admin/api/members", fa.handleAdminApiGetMembers).Methods(http.MethodGet)
	fa.router.HandleFunc("/admin/api/member/activate", fa.handleMemberActivate).Methods(http.MethodPut)
}
