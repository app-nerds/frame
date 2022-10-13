package webapp

import (
	"net/http"

	"github.com/app-nerds/frame/internal/baseviewmodel"
	"github.com/app-nerds/frame/internal/routepaths"
	webapp "github.com/app-nerds/frame/pkg/web-app"
	"github.com/gorilla/sessions"
)

func (wa *WebApp) handleAdminDashboard(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"AppName": wa.appName,
	}

	wa.RenderTemplate(w, "admin-dashboard.tmpl", data)
}

func (wa *WebApp) handleAdminLogin(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		session *sessions.Session
	)

	data := AdminLoginData{
		BaseViewModel: baseviewmodel.BaseViewModel{
			JavascriptIncludes: []webapp.JavascriptInclude{},
			AppName:            wa.appName,
			Stylesheets:        []string{},
		},
		Message: "",
	}

	/*
	 * Handle login submission
	 */
	if r.Method == http.MethodPost {
		_ = r.ParseForm()

		// First, is this a root user?
		if r.FormValue("userName") == wa.frameConfig.RootUserName && r.FormValue("password") == wa.frameConfig.RootUserPassword {
			if session, err = wa.adminSessionStore.Get(r, wa.adminSessionName); err != nil {
				wa.logger.WithError(err).Error("error geting session")
				http.Redirect(w, r, routepaths.UnexpectedErrorPath, http.StatusFound)
				return
			}

			session.Values["adminUserName"] = wa.frameConfig.RootUserName

			if err = wa.adminSessionStore.Save(r, w, session); err != nil {
				wa.logger.WithError(err).Error("error saving session")
				http.Redirect(w, r, routepaths.UnexpectedErrorPath, http.StatusFound)
				return
			}

			goTo := r.FormValue("referer")

			if goTo == "" {
				goTo = "/admin"
			}

			http.Redirect(w, r, goTo, http.StatusFound)
			return
		}

		// TODO: Add users from a database

		data.Message = "Invalid user name or password"
	}

	wa.RenderTemplate(w, "admin-login.tmpl", data)
}
