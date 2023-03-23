package frame

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

type InternalSiteAuthConfig struct {
	FrameStaticFS fs.FS
	Logger        *logrus.Entry
	SessionName   string
	SessionStore  sessions.Store
}

type SiteAuth struct {
	contentTemplateName   string
	frameStaticFS         fs.FS
	htmlPaths             []string
	layoutName            string
	logger                *logrus.Entry
	pathsExcludedFromAuth []string
	sessionName           string
	sessionStore          sessions.Store
}

/*
NewSiteAuth sets up authentication for the public site. This authentication makes use of an internal email/password
mechanism. It registers the following endpoints, each with baked in HTML:

- /member/login
- /member/account-pending
- /member/create-account
- /api/member/current
- /api/member/logout

All pages in SiteAuth have a few expectations.

1. They define a content area called "Title". This is so layouts can use "Title" to set the page title
2. These pages expect to use a master layout called "layout"
2. They define a content area called "content". They expect the primary layout to have a section for "content"
*/
func NewSiteAuth(internalConfig InternalSiteAuthConfig, siteAuthConfig SiteAuthConfig) *SiteAuth {
	result := &SiteAuth{
		contentTemplateName:   siteAuthConfig.ContentTemplateName,
		frameStaticFS:         internalConfig.FrameStaticFS,
		htmlPaths:             siteAuthConfig.HtmlPaths,
		layoutName:            siteAuthConfig.LayoutName,
		logger:                internalConfig.Logger,
		pathsExcludedFromAuth: siteAuthConfig.PathsExcludedFromAuth,
		sessionName:           internalConfig.SessionName,
		sessionStore:          internalConfig.SessionStore,
	}

	/*
	 * Make sure specific paths are excluded from auth
	 */

	// TODO: figure out admin auth
	result.pathsExcludedFromAuth = append(result.pathsExcludedFromAuth, "/static", "/admin-static", "/frame-static", SiteAuthAccountPendingPath, SiteAuthLoginPath,
		SiteAuthLogoutPath, MemberSignUpPath, UnexpectedErrorPath, "/admin")

	// These paths need to redirect to an HTML error or login page when the user is not authorized
	result.htmlPaths = append(result.htmlPaths, "/member/profile", "/member/profile/avatar")

	return result
}

func (sa *SiteAuth) RegisterStaticFrameAssetsRoute(router *mux.Router) {
	sa.logger.Info("registering static frame assets...")
	frameStaticFS := http.FileServer(http.FS(sa.frameStaticFS))
	router.PathPrefix("/frame-static/").Handler(frameStaticFS).Methods(http.MethodGet)
}

func (sa *SiteAuth) RegisterSiteAuthRoutes(router *mux.Router, webApp *WebApp, memberService *MemberService) {
	router.HandleFunc(SiteAuthLoginPath, sa.handleSiteAuthLogin(webApp, memberService)).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc(SiteAuthAccountPendingPath, sa.handleAccountPending(webApp)).Methods(http.MethodGet)

	sa.setupMiddleware(router)
}

func (sa *SiteAuth) setupMiddleware(router *mux.Router) {
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
			for _, excludedPath := range sa.pathsExcludedFromAuth {
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
			if session, err = sa.sessionStore.Get(r, sa.sessionName); err != nil {
				sa.logger.WithError(err).Error("error getting session information")
				http.Redirect(w, r, UnexpectedErrorPath, http.StatusFound)
				return
			}

			email, ok = session.Values["email"].(string)

			if !ok {
				sa.logger.WithFields(logrus.Fields{
					"ip":   RealIP(r),
					"path": r.URL.Path,
				}).Error("user is not authorized")

				sa.sendUnauthorizedResponse(w, r, sa.htmlPaths)
				return
			}

			if email == "" {
				sa.logger.WithFields(logrus.Fields{
					"ip":   RealIP(r),
					"path": r.URL.Path,
				}).Error("user is not authorized")

				sa.sendUnauthorizedResponse(w, r, sa.htmlPaths)
				return
			}

			status, ok := session.Values["status"].(string)
			if !ok {
				status = ""
			}

			if status != string(MemberActive) {
				sa.logger.WithFields(logrus.Fields{
					"ip":   RealIP(r),
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

	router.Use(middleware)
}

func (sa *SiteAuth) sendUnauthorizedResponse(w http.ResponseWriter, r *http.Request, htmlResponsePaths []string) {
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
