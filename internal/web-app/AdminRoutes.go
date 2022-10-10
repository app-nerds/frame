package webapp

import (
	"net/http"
)

func (wa *WebApp) handleAdminDashboard(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"appName": wa.appName,
	}

	wa.RenderTemplate(w, "admin-dashboard.tmpl", data)
}
