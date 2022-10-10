package siteauth

import (
	"errors"
	"net/http"

	"github.com/app-nerds/frame/internal/routepaths"
	webapp "github.com/app-nerds/frame/internal/web-app"
	"github.com/app-nerds/frame/pkg/framemember"
	"github.com/gorilla/sessions"
	"gorm.io/gorm"
)

func (sa *SiteAuth) handleSiteAuthLogin(webApp *webapp.WebApp, memberService *framemember.MemberService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err     error
			session *sessions.Session
			member  framemember.Member
		)

		data := map[string]interface{}{}

		for k, v := range sa.baseData {
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
			member, err = memberService.GetMemberByEmail(email, true)

			if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
				http.Redirect(w, r, routepaths.MemberSignUpPath, http.StatusFound)
				return
			}

			if err != nil {
				sa.logger.WithError(err).Error("error getting member information in site auth")
				http.Redirect(w, r, routepaths.UnexpectedErrorPath, http.StatusFound)
				return
			}

			/*
			 * If we have an existing member, but they aren't approved yet,
			 * redirect them.
			 */
			if member.Status.ID != framemember.MemberActiveID {
				http.Redirect(w, r, routepaths.SiteAuthAccountPendingPath, http.StatusFound)
				return
			}

			/*
			 * If they are approved, but deleted, redirect
			 */
			if member.DeletedAt.Valid {
				http.Redirect(w, r, routepaths.SiteAuthAccountPendingPath, http.StatusFound)
				return
			}

			/*
			 * If we have an approved member, but the password is invalid, let them know
			 */
			if !member.Password.IsSameAsPlaintextPassword(password) {
				data["errorMessage"] = "Invalid user name or password. Please try again."

				webApp.RenderTemplate(w, "login.tmpl", data)
				return
			}

			/*
			 * Otherwise, we are good to go!
			 */
			if session, err = sa.sessionStore.Get(r, sa.sessionName); err != nil {
				sa.logger.WithError(err).Error("error geting session")
				http.Redirect(w, r, routepaths.UnexpectedErrorPath, http.StatusFound)
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
				http.Redirect(w, r, routepaths.UnexpectedErrorPath, http.StatusFound)
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
