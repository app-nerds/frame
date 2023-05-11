package frame

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/app-nerds/kit/v6/datetime"
	"github.com/gorilla/mux"
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

func GetIntFromRequest(r *http.Request, name string) int {
	var (
		err         error
		valueString string
		value       int
	)

	valueString = r.FormValue(name)

	if value, err = strconv.Atoi(valueString); err != nil {
		vars := mux.Vars(r)
		valueString = vars[name]

		if value, err = strconv.Atoi(valueString); err != nil {
			value = 0
		}
	}

	return value
}

func GetFloatFromRequest(r *http.Request, name string) float64 {
	var (
		err         error
		valueString string
		value       float64
	)

	valueString = r.FormValue(name)

	if value, err = strconv.ParseFloat(valueString, 64); err != nil {
		vars := mux.Vars(r)
		valueString = vars[name]

		if value, err = strconv.ParseFloat(valueString, 64); err != nil {
			value = 0
		}
	}

	return value
}

func GetStringFromRequest(r *http.Request, name string) string {
	value := r.FormValue(name)

	if value == "" {
		vars := mux.Vars(r)
		value = vars[name]
	}

	return value
}

func GetTimeFromRequest(r *http.Request, name string) time.Time {
	var (
		err         error
		valueString string
		value       time.Time
	)

	parser := datetime.DateTimeParser{}
	valueString = r.FormValue(name)

	if value, err = parser.Parse(valueString); err != nil {
		vars := mux.Vars(r)
		valueString = vars[name]

		if value, err = parser.Parse(valueString); err != nil {
			return time.Now().UTC()
		}
	}

	return value
}

func GetPageFromRequest(r *http.Request) int {
	page := GetIntFromRequest(r, "page")

	if page < 1 {
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
