package siteauth

import (
	"net/http"

	webapp "github.com/app-nerds/frame/internal/web-app"
)

func handleAuthAccountPending(webApp *webapp.WebApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		webApp.RenderTemplate(w, "account-pending.tmpl", map[string]interface{}{})
	}
}
