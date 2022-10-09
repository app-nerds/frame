package frame

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

func (fa *FrameApplication) setupMiddleware(pathsExcludedFromAuth, htmlPaths []string) {
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				err     error
				session *sessions.Session
				ok      bool

				email string
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
			if session, err = fa.sessionStore.Get(r, fa.sessionName); err != nil {
				fa.Logger.WithError(err).Error("error getting session information")
				http.Redirect(w, r, UnexpectedErrorPath, http.StatusFound)
				return
			}

			email, ok = session.Values["email"].(string)

			if !ok {
				fa.Logger.WithFields(logrus.Fields{
					"ip":   realIP(r),
					"path": r.URL.Path,
				}).Error("user is not authorized")

				fa.sendUnauthorizedResponse(w, r, htmlPaths)
				return
			}

			if email == "" {
				fa.Logger.WithFields(logrus.Fields{
					"ip":   realIP(r),
					"path": r.URL.Path,
				}).Error("user is not authorized")

				fa.sendUnauthorizedResponse(w, r, htmlPaths)
				return
			}

			status, ok := session.Values["status"].(string)
			if !ok {
				status = ""
			}

			if status != string(MemberActive) {
				fa.Logger.WithFields(logrus.Fields{
					"ip":   realIP(r),
					"path": r.URL.Path,
				}).Error("user has an account but it is not yet approved")

				http.Redirect(w, r, SiteAuthAccountPendingPath, http.StatusFound)
				return
			}

			memberID, ok := session.Values["memberID"].(uint)
			if !ok {
				memberID = 0
			}

			firstName, ok := session.Values["firstName"].(string)
			if !ok {
				firstName = ""
			}

			lastName, ok := session.Values["lastName"].(string)
			if !ok {
				lastName = ""
			}

			avatarURL, ok := session.Values["avatarURL"].(string)
			if !ok {
				avatarURL = ""
			}

			ctx := context.WithValue(r.Context(), "firstName", firstName)
			ctx = context.WithValue(ctx, "lastName", lastName)
			ctx = context.WithValue(ctx, "email", email)
			ctx = context.WithValue(ctx, "avatarURL", avatarURL)
			ctx = context.WithValue(ctx, "memberID", memberID)
			ctx = context.WithValue(ctx, "status", status)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	fa.router.Use(middleware)
}

func (fa *FrameApplication) sendUnauthorizedResponse(w http.ResponseWriter, r *http.Request, htmlResponsePaths []string) {
	for _, path := range htmlResponsePaths {
		if strings.HasPrefix(r.URL.Path, path) {
			http.Redirect(w, r, fmt.Sprintf("%s?referer=%s", SiteAuthLoginPath, r.URL.Path), http.StatusFound)
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
