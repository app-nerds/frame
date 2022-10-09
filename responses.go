package frame

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GenericErrorResponse struct {
	Code    string `json:"code"`
	Detail  string `json:"detail"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type GenericSuccessResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func (fa *FrameApplication) CreateGenericErrorResponse(message, detail, code string) GenericErrorResponse {
	return GenericErrorResponse{
		Code:    code,
		Detail:  detail,
		Message: message,
		Success: false,
	}
}

func (fa *FrameApplication) CreateGenericSuccessResponse(message string) GenericSuccessResponse {
	return GenericSuccessResponse{
		Message: message,
		Success: true,
	}
}

/*
ReadJSONBody reads the body content from an http.Request as JSON data into
dest.
*/
func (fa *FrameApplication) ReadJSONBody(r *http.Request, dest interface{}) error {
	var (
		err error
		b   []byte
	)

	if b, err = io.ReadAll(r.Body); err != nil {
		return fmt.Errorf("error reading request body: %w", err)
	}

	if err = json.Unmarshal(b, &dest); err != nil {
		return fmt.Errorf("error unmarshaling body to destination: %w", err)
	}

	return nil
}

/*
WriteJSON writes JSON content to the response writer.
*/
func (fa *FrameApplication) WriteJSON(w http.ResponseWriter, status int, value interface{}) {
	var (
		err error
		b   []byte
	)

	w.Header().Set("Content-Type", "application/json")

	if b, err = json.Marshal(value); err != nil {
		fa.Logger.WithError(err).Error("error marshaling value for writing")

		b, _ = json.Marshal(struct {
			Message    string `json:"message"`
			Suggestion string `json:"suggestion"`
		}{
			Message:    "Error marshaling value for writing",
			Suggestion: "See error log for more information",
		})

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "%s", string(b))
		return
	}

	if status > 299 {
		w.WriteHeader(status)
	}

	_, _ = fmt.Fprintf(w, "%s", string(b))
}

/*
WriteString writes string content to the response writer.
*/
func (fa *FrameApplication) WriteString(w http.ResponseWriter, status int, value string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	_, _ = fmt.Fprintf(w, "%s", value)
}
