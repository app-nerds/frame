package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/app-nerds/frame/internal/routepaths"
	"github.com/app-nerds/frame/pkg/config"
	"github.com/app-nerds/frame/pkg/httputils"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

func AdminAuthMiddleware(logger *logrus.Entry, config *config.Config, sessionStore sessions.Store) mux.MiddlewareFunc {
	pathsExcludedFromAuth := []string{
		"/admin/login",
		"/frame-static/",
		"/admin-static/",
	}

	htmlPaths := []string{
		"/admin",
		"/members/manage",
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				err     error
				session *sessions.Session
				ok      bool

				adminUserName string
			)

			/*
			 * If this path is excluded from auth, just keep going
			 */
			for _, excludedPath := range pathsExcludedFromAuth {
				if excludedPath == "/" && r.URL.Path == "/" {
					next.ServeHTTP(w, r)
					return
				}

				if strings.HasPrefix(r.URL.Path, excludedPath) && excludedPath != "/" {
					next.ServeHTTP(w, r)
					return
				}
			}

			/*
			 * If not, let's verify we have a cookie
			 */
			if session, err = sessionStore.Get(r, config.AdminSessionName); err != nil {
				logger.WithError(err).Error("error getting admin session information")
				http.Redirect(w, r, routepaths.UnexpectedErrorPath, http.StatusFound)
				return
			}

			adminUserName, ok = session.Values["adminUserName"].(string)

			if !ok {
				logger.WithFields(logrus.Fields{
					"ip":   httputils.RealIP(r),
					"path": r.URL.Path,
				}).Error("user is not authorized")

				adminMiddlewareSendUnauthorizedResponse(w, r, htmlPaths)
				return
			}

			if adminUserName == "" {
				logger.WithFields(logrus.Fields{
					"ip":   httputils.RealIP(r),
					"path": r.URL.Path,
				}).Error("user is not authorized")

				adminMiddlewareSendUnauthorizedResponse(w, r, htmlPaths)
				return
			}

			ctx := context.WithValue(r.Context(), "adminUserName", adminUserName)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func adminMiddlewareSendUnauthorizedResponse(w http.ResponseWriter, r *http.Request, htmlResponsePaths []string) {
	for _, path := range htmlResponsePaths {
		if strings.HasPrefix(r.URL.Path, path) {
			http.Redirect(w, r, fmt.Sprintf("%s?referer=%s", routepaths.AdminLoginPath, r.URL.Path), http.StatusFound)
			return
		}
	}

	result := map[string]interface{}{
		"success": false,
		"error":   "User unauthorized",
		"status":  http.StatusUnauthorized,
	}

	b, _ := json.Marshal(result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = fmt.Fprint(w, string(b))
}
