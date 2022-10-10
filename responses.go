package frame

import (
	"net/http"

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
