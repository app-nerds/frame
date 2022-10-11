package frame

import (
	"net/http"

	"github.com/app-nerds/frame/internal/routepaths"
	"github.com/app-nerds/frame/pkg/httputils"
)

func (fa *FrameApplication) CreateGenericErrorResponse(message, detail, code string) httputils.GenericErrorResponse {
	return httputils.GenericErrorResponse{
		Code:    code,
		Detail:  detail,
		Message: message,
		Success: false,
	}
}

func (fa *FrameApplication) CreateGenericSuccessResponse(message string) httputils.GenericSuccessResponse {
	return httputils.GenericSuccessResponse{
		Message: message,
		Success: true,
	}
}

/*
ReadJSONBody reads the body content from an http.Request as JSON data into
dest.
*/
func (fa *FrameApplication) ReadJSONBody(r *http.Request, dest interface{}) error {
	return httputils.ReadJSONBody(r, dest)
}

/*
WriteJSON writes JSON content to the response writer.
*/
func (fa *FrameApplication) WriteJSON(w http.ResponseWriter, status int, value interface{}) {
	httputils.WriteJSON(w, status, value)
}

/*
WriteString writes string content to the response writer.
*/
func (fa *FrameApplication) WriteString(w http.ResponseWriter, status int, value string) {
	httputils.WriteString(w, status, value)
}

/*
UnexpectedError redirects the user to a page for unexpected errors. This is configured
when calling AddWebApp
*/
func (fa *FrameApplication) UnexpectedError(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, routepaths.UnexpectedErrorPath, http.StatusFound)
}
