package frame

import (
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
	"gorm.io/gorm"
)

/*
SiteAuth sets up authentication for the public site. This authentication makes use of an internal email/password
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
func (fa *FrameApplication) SetupSiteAuth(layoutName, contentTemplateName string, baseData map[string]interface{}, pathsExcludedFromAuth, htmlPaths []string) *FrameApplication {
	if fa.sessionStore == nil {
		fa.Logger.Fatal("Please setup a session storage before calling SetupSiteAuth()")
	}

	/*
	 * Make sure specific paths are excluded from auth
	 */
	pathsExcludedFromAuth = append(pathsExcludedFromAuth, "/static", SiteAuthAccountPendingPath, SiteAuthLoginPath, SiteAuthLogoutPath, SiteAuthMemberSignUpPath, UnexpectedErrorPath)

	fa.router.HandleFunc(SiteAuthAccountPendingPath, fa.handleAuthAccountPending).Methods(http.MethodGet)
	fa.router.HandleFunc(SiteAuthLoginPath, fa.handleSessionBasicLogin(baseData)).Methods(http.MethodGet, http.MethodPost)
	fa.router.HandleFunc(SiteAuthMemberSignUpPath, fa.handleMemberSignup).Methods(http.MethodGet, http.MethodPost)
	fa.router.HandleFunc(MemberApiCurrentMember, fa.handleMemberCurrent).Methods(http.MethodGet)
	fa.router.HandleFunc(MemberApiLogOut, fa.handleMemberLogout).Methods(http.MethodGet)

	fa.setupMiddleware(pathsExcludedFromAuth, htmlPaths)
	return fa
}

func (fa *FrameApplication) handleAuthAccountPending(w http.ResponseWriter, r *http.Request) {
	fa.RenderTemplate(w, "account-pending.tmpl", map[string]interface{}{})
}

func (fa *FrameApplication) handleSessionBasicLogin(baseData map[string]interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			 * If this member doesn't exist yet, send them to the sign-up page
			 */
			member, err = fa.MemberService.GetMemberByEmail(email, true)

			if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
				http.Redirect(w, r, SiteAuthMemberSignUpPath, http.StatusFound)
				return
			}

			if err != nil {
				fa.Logger.WithError(err).Error("error getting member information in site auth")
				http.Redirect(w, r, UnexpectedErrorPath, http.StatusFound)
				return
			}

			/*
			 * If we have an existing member, but they aren't approved yet,
			 * redirect them.
			 */
			if !member.Approved {
				http.Redirect(w, r, SiteAuthAccountPendingPath, http.StatusFound)
				return
			}

			/*
			 * If they are approved, but deleted, redirect
			 */
			if member.DeletedAt.Valid {
				http.Redirect(w, r, SiteAuthAccountPendingPath, http.StatusFound)
				return
			}

			/*
			 * If we have an approved member, but the password is invalid, let them know
			 */
			if !member.Password.IsSameAsPlaintextPassword(password) {
				data["errorMessage"] = "Invalid user name or password. Please try again."

				fa.RenderTemplate(w, "login.tmpl", data)
				return
			}

			/*
			 * Otherwise, we are good to go!
			 */
			if session, err = fa.sessionStore.Get(r, fa.sessionName); err != nil {
				fa.Logger.WithError(err).Error("error geting session")
				http.Redirect(w, r, UnexpectedErrorPath, http.StatusFound)
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
				http.Redirect(w, r, UnexpectedErrorPath, http.StatusFound)
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
	}
}
