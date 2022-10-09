package frame

import "net/http"

/*
GET /errors/unexpected
*/
func (fa *FrameApplication) handleUnexpectedError(w http.ResponseWriter, r *http.Request) {
	fa.RenderTemplate(w, "unexpected-error.tmpl", nil)
}
