package frame

import (
	"net/http"
)

func handleAuthAccountPending(webApp *WebApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		webApp.RenderTemplate(w, "account-pending.tmpl", map[string]interface{}{})
	}
}
