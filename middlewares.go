package frame

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const (
	AllowAllOrigins           string = "*"
	AllowAllMethods           string = "POST, GET, OPTIONS, PUT, DELETE"
	AllowAllHeaders           string = "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"
	AllowHeaderAccept         string = "Accept"
	AllowHeaderContentType    string = "Content-Type"
	AllowHeaderContentLength  string = "Content-Length"
	AllowHeaderAcceptEncoding string = "Accept-Encoding"
	AllowHeaderCSRF           string = "X-CSRF-Token"
	AllowHeaderAuthorization  string = "Authorization"
)

type accessControl struct {
	handler      http.Handler
	allowOrigin  string
	allowMethods string
	allowHeaders string
}

func (ac *accessControl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", ac.allowOrigin)
	w.Header().Set("Access-Control-Allow-Methods", ac.allowMethods)
	w.Header().Set("Access-Control-Allow-Headers", ac.allowHeaders)

	ac.handler.ServeHTTP(w, r)
}

/*
accessControlMiddleware wraps an HTTP mux with a middleware that sets
headers for access control and allowed headers.
*/
func accessControlMiddleware(allowOrigin, allowMethods, allowHeaders string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler := &accessControl{
				handler:      next,
				allowOrigin:  allowOrigin,
				allowMethods: allowMethods,
				allowHeaders: allowHeaders,
			}

			handler.ServeHTTP(w, r)
		})
	}
}

/*
Allow verifies if the caller method matches the provided method.

If the caller's method does not match what is allowed, the string
"method not allowed" is written back to the caller.
*/
func allow(next http.HandlerFunc, allowedMethod string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.ToLower(r.Method) != strings.ToLower(allowedMethod) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = fmt.Fprintf(w, "%s", "method not allowed")

			return
		}

		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}
