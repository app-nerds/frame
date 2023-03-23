package frame

import (
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
	"gorm.io/gorm"
)

func (sa *SiteAuth) handleSiteAuthLogin(webApp *WebApp, memberService *MemberService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err     error
			session *sessions.Session
			member  Member
		)

		data := struct {
			Email        string
			ErrorMessage string
			Referer      string
			Stylesheets  []string
		}{
			Stylesheets: []string{
				"/frame-static/css/frame-page-styles.css",
			},
		}

		if r.Method == http.MethodGet {
			data.Referer = r.URL.Query().Get("referer")
		}

		if r.Method == http.MethodPost {
			r.ParseForm()

			data.Email = r.FormValue("email")
			data.Referer = r.Form.Get("referer")
			email := r.Form.Get("email")
			password := r.Form.Get("password")

			/*
			 * If this member doesn't exist yet, tell them they can make one.
			 */
			member, err = memberService.GetMemberByEmail(email, true)

			if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
				data.ErrorMessage = "Invalid user name or password. Please try again."
				webApp.RenderTemplate(w, "login.tmpl", data)
				return
			}

			if err != nil {
				sa.logger.WithError(err).Error("error getting member information in site auth")
				http.Redirect(w, r, UnexpectedErrorPath, http.StatusFound)
				return
			}

			/*
			 * If we have an existing member, but they aren't approved yet,
			 * redirect them.
			 */
			if member.Status.ID != MemberActiveID {
				http.Redirect(w, r, SiteAuthAccountPendingPath, http.StatusFound)
				return
			}

			/*
			 * If they are approved, but deleted, redirect
			 */
			if member.DeletedAt != nil {
				http.Redirect(w, r, SiteAuthAccountPendingPath, http.StatusFound)
				return
			}

			/*
			 * If we have an approved member, but the password is invalid, let them know
			 */
			if !member.Password.IsSameAsPlaintextPassword(password) {
				data.ErrorMessage = "Invalid user name or password. Please try again."
				webApp.RenderTemplate(w, "login.tmpl", data)
				return
			}

			/*
			 * Otherwise, we are good to go!
			 */
			if session, err = sa.sessionStore.Get(r, sa.sessionName); err != nil {
				sa.logger.WithError(err).Error("error geting session")
				http.Redirect(w, r, UnexpectedErrorPath, http.StatusFound)
				return
			}

			session.Values["memberID"] = member.ID
			session.Values["email"] = member.Email
			session.Values["firstName"] = member.FirstName
			session.Values["lastName"] = member.LastName
			session.Values["avatarURL"] = member.AvatarURL
			session.Values["status"] = string(member.Status.Status)

			if err = sa.sessionStore.Save(r, w, session); err != nil {
				sa.logger.WithError(err).Error("error saving session")
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

		webApp.RenderTemplate(w, "login.tmpl", data)
	}
}

func (sa *SiteAuth) handleAccountPending(webApp *WebApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Stylesheets []string
		}{
			Stylesheets: []string{
				"/frame-static/css/frame-page-styles.css",
			},
		}

		webApp.RenderTemplate(w, "account-pending.tmpl", data)
	}
}
