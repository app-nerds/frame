package frame

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

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

func (fa *FrameApplication) SetupExternalAuth(pathsExcludedFromAuth, htmlPaths []string) *FrameApplication {
	if fa.sessionStore == nil {
		fa.Logger.Fatal("Please setup a session storage before calling WithExternalAuth()")
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
			http.Redirect(w, r, fa.unexpectedErrorPath, http.StatusTemporaryRedirect)
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
				http.Redirect(w, r, fa.unexpectedErrorPath, http.StatusTemporaryRedirect)
				return
			}

			http.Redirect(w, r, fa.accountAwaitingApprovalPath, http.StatusTemporaryRedirect)
			return
		}

		if err != nil {
			fa.Logger.WithError(err).Error("error getting member information in external auth")
			http.Redirect(w, r, fa.unexpectedErrorPath, http.StatusTemporaryRedirect)
			return
		}

		/*
		 * If we have an existing member, but they aren't approved yet,
		 * redirect them.
		 */
		if !member.Approved {
			http.Redirect(w, r, fa.accountAwaitingApprovalPath, http.StatusTemporaryRedirect)
			return
		}

		/*
		 * Otherwise, we are good to go!
		 */
		if session, err = fa.sessionStore.Get(r, fa.sessionName); err != nil {
			fa.Logger.WithError(err).Error("error geting session")
			http.Redirect(w, r, fa.unexpectedErrorPath, http.StatusTemporaryRedirect)
			return
		}

		session.Values["email"] = user.Email
		session.Values["firstName"] = user.FirstName
		session.Values["lastName"] = user.LastName
		session.Values["avatarURL"] = user.AvatarURL
		session.Values["approved"] = member.Approved

		if err = fa.sessionStore.Save(r, w, session); err != nil {
			fa.Logger.WithError(err).Error("error saving session")
			http.Redirect(w, r, fa.unexpectedErrorPath, http.StatusTemporaryRedirect)
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
				http.Redirect(w, r, fa.unexpectedErrorPath, http.StatusTemporaryRedirect)
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

				http.Redirect(w, r, fa.accountAwaitingApprovalPath, http.StatusTemporaryRedirect)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	fa.router.Use(middleware)
}

func (fa *FrameApplication) sendUnauthorizedResponse(w http.ResponseWriter, r *http.Request, htmlResponsePaths []string) {
	for _, path := range htmlResponsePaths {
		if strings.HasPrefix(r.URL.Path, path) {
			http.Redirect(w, r, fa.unauthorizedPath, http.StatusTemporaryRedirect)
			return
		}
	}

	result := map[string]interface{}{
		"success": false,
		"error":   "User unauthorized",
	}

	b, _ := json.Marshal(result)

	w.Header().Set("Content-Type", "application/json")
	_, _ = fmt.Fprint(w, b)
}
