package frame

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/app-nerds/kit/v6/passwords"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (fa *FrameApplication) WithGoogleAuth(scopes ...string) *FrameApplication {
	fa.externalAuths = append(fa.externalAuths, google.New(fa.Config.GoogleClientID, fa.Config.GoogleClientSecret, fa.Config.GoogleRedirectURI, scopes...))
	return fa
}

/*
SiteAuth sets up authentication for the public site. This authentication makes use of an internal email/password
mechanism. It registers the following endpoints, each with baked in HTML:

- /member/login
- /member/logout
- /member/account-pending

All pages in SiteAuth have a few expectations.

1. They define a content area called "Title". This is so layouts can use "Title" to set the page title
2. They make use only of the PrimaryLayoutName setting.
2. They define a content area called "content". They expect the primary layout to have a section for "content"
*/
func (fa *FrameApplication) SetupSiteAuth(layoutName, contentTemplateName string, baseData map[string]interface{}, pathsExcludedFromAuth, htmlPaths []string) *FrameApplication {
	if fa.sessionStore == nil {
		fa.Logger.Fatal("Please setup a session storage before calling SetupSiteAuth()")
	}

	/*
	 * Make sure specific paths are excluded from auth
	 */
	pathsExcludedFromAuth = append(pathsExcludedFromAuth, "/static", SiteAuthAccountPendingPath, SiteAuthLoginPath, SiteAuthLogoutPath, UnexpectedErrorPath)

	fa.router.HandleFunc(SiteAuthAccountPendingPath, func(w http.ResponseWriter, r *http.Request) {
		fa.RenderTemplate(w, "account-pending.tmpl", map[string]interface{}{})
	})

	fa.router.HandleFunc(SiteAuthLoginPath, func(w http.ResponseWriter, r *http.Request) {
		var (
			err     error
			session *sessions.Session
			member  Member
		)

		data := map[string]interface{}{}

		for k, v := range baseData {
			data[k] = v
		}

		if r.Method == http.MethodGet {
			data["Referer"] = r.URL.Query().Get("referer")
		}

		if r.Method == http.MethodPost {
			r.ParseForm()

			data["Referer"] = r.Form.Get("referer")
			email := r.Form.Get("email")
			password := r.Form.Get("password")

			/*
			 * If this member doesn't exist yet, create them as an unapproved member
			 */
			member, err = fa.MemberService.GetMemberByEmail(email)

			if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
				member = Member{
					Approved:  false,
					AvatarURL: "",
					Email:     email,
					Password:  passwords.HashedPasswordString(password),
				}

				if err = fa.MemberService.CreateMember(&member); err != nil {
					fa.Logger.WithError(err).Error("error creating new member")
					http.Redirect(w, r, UnexpectedErrorPath, http.StatusTemporaryRedirect)
					return
				}

				http.Redirect(w, r, SiteAuthAccountPendingPath, http.StatusTemporaryRedirect)
				return
			}

			if err != nil {
				fa.Logger.WithError(err).Error("error getting member information in site auth")
				http.Redirect(w, r, UnexpectedErrorPath, http.StatusTemporaryRedirect)
				return
			}

			/*
			 * If we have an existing member, but they aren't approved yet,
			 * redirect them.
			 */
			if !member.Approved {
				http.Redirect(w, r, SiteAuthAccountPendingPath, http.StatusTemporaryRedirect)
				return
			}

			/*
			 * If we have an approved member, but the password is invalid, let them know
			 */
			if !member.Password.IsSameAsPlaintextPassword(password) {
				data["ErrorMessage"] = "Invalid user name or password. Please try again."

				fa.RenderTemplate(w, "login.tmpl", data)
				return
			}

			/*
			 * Otherwise, we are good to go!
			 */
			if session, err = fa.sessionStore.Get(r, fa.sessionName); err != nil {
				fa.Logger.WithError(err).Error("error geting session")
				http.Redirect(w, r, UnexpectedErrorPath, http.StatusTemporaryRedirect)
				return
			}

			session.Values["memberID"] = member.ID
			session.Values["email"] = member.Email
			session.Values["firstName"] = member.FirstName
			session.Values["lastName"] = member.LastName
			session.Values["avatarURL"] = member.AvatarURL
			session.Values["approved"] = member.Approved

			if err = fa.sessionStore.Save(r, w, session); err != nil {
				fa.Logger.WithError(err).Error("error saving session")
				http.Redirect(w, r, UnexpectedErrorPath, http.StatusTemporaryRedirect)
				return
			}

			goTo := r.Form.Get("referer")

			if goTo == "" {
				goTo = "/"
			}

			http.Redirect(w, r, goTo, http.StatusFound)
			return
		}

		fa.RenderTemplate(w, "login.tmpl", data)
	}).Methods(http.MethodGet, http.MethodPost)

	fa.router.HandleFunc(MemberApiCurrentMember, func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		member := map[string]interface{}{
			"id":        ctx.Value("memberID"),
			"email":     ctx.Value("email"),
			"firstName": ctx.Value("firstName"),
			"lastName":  ctx.Value("lastName"),
			"avatarURL": ctx.Value("avatarURL"),
		}

		fa.WriteJSON(w, http.StatusOK, member)
	})

	fa.router.HandleFunc(MemberApiLogOut, func(w http.ResponseWriter, r *http.Request) {
		var (
			err     error
			session *sessions.Session
		)

		if session, err = fa.sessionStore.Get(r, fa.sessionName); err != nil {
			fa.Logger.WithError(err).Error("error getting session information")
			fa.WriteJSON(w, http.StatusInternalServerError, fa.CreateGenericErrorResponse("error getting session information", err.Error(), ""))
			return
		}

		session.Options.MaxAge = -1

		if err = fa.sessionStore.Save(r, w, session); err != nil {
			fa.Logger.WithError(err).Error("error deleting session")
			fa.WriteJSON(w, http.StatusInternalServerError, fa.CreateGenericErrorResponse("error deleting session", err.Error(), ""))
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	})

	fa.setupMiddleware(pathsExcludedFromAuth, htmlPaths)
	return fa
}

func (fa *FrameApplication) SetupExternalAuth(pathsExcludedFromAuth, htmlPaths []string) *FrameApplication {
	if fa.sessionStore == nil {
		fa.Logger.Fatal("Please setup a session storage before calling SetupExternalAuth()")
	}

	gothic.Store = fa.sessionStore

	goth.UseProviders(
		fa.externalAuths...,
	)

	fa.router.HandleFunc("/auth/{provider}/callback", func(w http.ResponseWriter, r *http.Request) {
		var (
			err     error
			user    goth.User
			session *sessions.Session
			member  Member
		)

		user, err = gothic.CompleteUserAuth(w, r)

		if err != nil {
			http.Redirect(w, r, UnexpectedErrorPath, http.StatusTemporaryRedirect)
			return
		}

		/*
		 * If this member doesn't exist yet, create them as an unapproved member
		 */
		member, err = fa.MemberService.GetMemberByEmail(user.Email)

		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			member = Member{
				Approved:   false,
				AvatarURL:  user.AvatarURL,
				Email:      user.Email,
				ExternalID: user.UserID,
				FirstName:  user.FirstName,
				LastName:   user.LastName,
			}

			if err = fa.MemberService.CreateMember(&member); err != nil {
				fa.Logger.WithError(err).Error("error creating new member")
				http.Redirect(w, r, UnexpectedErrorPath, http.StatusTemporaryRedirect)
				return
			}

			http.Redirect(w, r, SiteAuthAccountPendingPath, http.StatusTemporaryRedirect)
			return
		}

		if err != nil {
			fa.Logger.WithError(err).Error("error getting member information in external auth")
			http.Redirect(w, r, UnexpectedErrorPath, http.StatusTemporaryRedirect)
			return
		}

		/*
		 * If we have an existing member, but they aren't approved yet,
		 * redirect them.
		 */
		if !member.Approved {
			http.Redirect(w, r, SiteAuthAccountPendingPath, http.StatusTemporaryRedirect)
			return
		}

		/*
		 * Otherwise, we are good to go!
		 */
		if session, err = fa.sessionStore.Get(r, fa.sessionName); err != nil {
			fa.Logger.WithError(err).Error("error geting session")
			http.Redirect(w, r, UnexpectedErrorPath, http.StatusTemporaryRedirect)
			return
		}

		session.Values["email"] = user.Email
		session.Values["firstName"] = user.FirstName
		session.Values["lastName"] = user.LastName
		session.Values["avatarURL"] = user.AvatarURL
		session.Values["approved"] = member.Approved

		if err = fa.sessionStore.Save(r, w, session); err != nil {
			fa.Logger.WithError(err).Error("error saving session")
			http.Redirect(w, r, UnexpectedErrorPath, http.StatusTemporaryRedirect)
			return
		}

		if fa.OnAuthSuccess != nil {
			fa.OnAuthSuccess(w, r, member)
		}
	})

	fa.router.HandleFunc("/auth/{provider}", func(w http.ResponseWriter, r *http.Request) {
		gothic.BeginAuthHandler(w, r)
	})

	fa.setupMiddleware(pathsExcludedFromAuth, htmlPaths)
	return fa
}

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
				http.Redirect(w, r, UnexpectedErrorPath, http.StatusTemporaryRedirect)
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

			approved, _ := session.Values["approved"].(bool)

			if !approved {
				fa.Logger.WithFields(logrus.Fields{
					"ip":   realIP(r),
					"path": r.URL.Path,
				}).Error("user has an account but it is not yet approved")

				http.Redirect(w, r, SiteAuthAccountPendingPath, http.StatusTemporaryRedirect)
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

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	fa.router.Use(middleware)
}

func (fa *FrameApplication) sendUnauthorizedResponse(w http.ResponseWriter, r *http.Request, htmlResponsePaths []string) {
	for _, path := range htmlResponsePaths {
		if strings.HasPrefix(r.URL.Path, path) {
			http.Redirect(w, r, fmt.Sprintf("%s?referer=%s", SiteAuthLoginPath, r.URL.Path), http.StatusTemporaryRedirect)
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
