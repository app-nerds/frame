package frame

import "net/http"

/*
GET /errors/unexpected
*/
func (wa *WebApp) handleUnexpectedError(w http.ResponseWriter, r *http.Request) {
	wa.RenderTemplate(w, "unexpected-error.tmpl", nil)
}
