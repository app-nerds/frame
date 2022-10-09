package frame

import (
	"net/http"
)

func (fa *FrameApplication) handleAdminDashboard(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"appName": fa.appName,
	}

	fa.RenderTemplate(w, "admin-dashboard.tmpl", data)
}
