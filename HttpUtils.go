package frame

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
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

func CreateGenericErrorResponse(message, detail, code string) GenericErrorResponse {
	return GenericErrorResponse{
		Code:    code,
		Detail:  detail,
		Message: message,
		Success: false,
	}
}

func CreateGenericSuccessResponse(message string) GenericSuccessResponse {
	return GenericSuccessResponse{
		Message: message,
		Success: true,
	}
}

func GetPageFromRequest(r *http.Request) int {
	var (
		err        error
		pageString string
		page       int
	)

	if r.Method == http.MethodGet {
		pageString = r.URL.Query().Get("page")
	}

	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		pageString = r.FormValue("page")
	}

	if page, err = strconv.Atoi(pageString); err != nil {
		page = 1
	}

	return page
}

/*
ReadJSONBody reads the body content from an http.Request as JSON data into
dest.
*/
func ReadJSONBody(r *http.Request, dest interface{}) error {
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
RealIP attempts to return the IP address of the caller. The result
will default to the RemoteAddr from http.Request. It will also
check the request headers for an "X-Forwarded-For" value and use
that. This is useful for when requests come through proxies or other
non-direct means.
*/
func RealIP(r *http.Request) string {
	result := r.RemoteAddr
	xForwardedFor := r.Header.Get("X-Forwarded-For")

	if xForwardedFor != "" {
		result = xForwardedFor
	}

	return result
}

/*
WriteJSON writes JSON content to the response writer.
*/
func WriteJSON(w http.ResponseWriter, status int, value interface{}) {
	var (
		err error
		b   []byte
	)

	w.Header().Set("Content-Type", "application/json")

	if b, err = json.Marshal(value); err != nil {
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
func WriteString(w http.ResponseWriter, status int, value string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	_, _ = fmt.Fprintf(w, "%s", value)
}
