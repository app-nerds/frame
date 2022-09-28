package frame

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
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

func realIP(r *http.Request) string {
	result := r.RemoteAddr
	xForwardedFor := r.Header.Get("X-Forwarded-For")

	if xForwardedFor != "" {
		result = xForwardedFor
	}

	return result
}

type statusRecorder struct {
	http.ResponseWriter
	Status int
}

func (sr *statusRecorder) Header() http.Header {
	return sr.ResponseWriter.Header()
}

func (sr *statusRecorder) Write(b []byte) (int, error) {
	return sr.ResponseWriter.Write(b)
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.Status = code
	sr.ResponseWriter.WriteHeader(code)
}

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

type requestLogger struct {
	handler http.Handler
	logger  *logrus.Entry
}

func (m *requestLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	recorder := &statusRecorder{
		ResponseWriter: w,
		Status:         http.StatusOK,
	}

	startTime := time.Now()
	ip := realIP(r)

	m.handler.ServeHTTP(recorder, r)
	diff := time.Since(startTime)

	m.logger.WithFields(logrus.Fields{
		"ip":            ip,
		"method":        r.Method,
		"status":        recorder.Status,
		"executionTime": diff,
		"queryParams":   r.URL.RawQuery,
	}).Info(r.URL.Path)
}

/*
RequestLogger returns a middleware for logging all requests.

Example:

  mux := nerdweb.NewServeMux()
  mux.HandleFunc("/endpoint", handler)

  mux.Use(middlewares.RequestLogger(logger))
*/
func requestLoggerMiddleware(logger *logrus.Entry) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler := &requestLogger{
				handler: next,
				logger:  logger,
			}

			handler.ServeHTTP(w, r)
		})
	}
}
