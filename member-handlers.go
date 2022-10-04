package frame

import (
	"net/http"

	"github.com/app-nerds/kit/v6/passwords"
	"github.com/gorilla/sessions"
)

/*
/api/member/current
*/
func (fa *FrameApplication) handleMemberCurrent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	member := map[string]interface{}{
		"id":        ctx.Value("memberID"),
		"email":     ctx.Value("email"),
		"firstName": ctx.Value("firstName"),
		"lastName":  ctx.Value("lastName"),
		"avatarURL": ctx.Value("avatarURL"),
	}

	fa.WriteJSON(w, http.StatusOK, member)
}

/*
/api/member/logout
*/
func (fa *FrameApplication) handleMemberLogout(w http.ResponseWriter, r *http.Request) {
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
}

/*
/member/create-account
*/
func (fa *FrameApplication) handleMemberSignup(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		member Member
	)

	data := map[string]interface{}{
		"errorMessage": "",
		"user": map[string]interface{}{
			"firstName": "",
			"lastName":  "",
			"email":     "",
		},
	}

	render := func() {
		fa.RenderTemplate(w, "sign-up.tmpl", data)
	}

	/*
	 * If we are just landing on the page, display it
	 */
	if r.Method == http.MethodGet {
		render()
		return
	}

	/*
	 * If we are posted here, sign the user up. Make sure we don't already have
	 * an existing user with this email address. If we do, let them know.
	 */
	_ = r.ParseForm()

	firstName := r.Form.Get("firstName")
	lastName := r.Form.Get("lastName")
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	reenterPassword := r.Form.Get("reenterPassword")

	member, err = fa.MemberService.GetMemberByEmail(email, true)

	// We already have a member with this email address
	if err == nil {
		data["user"] = map[string]interface{}{
			"firstName": firstName,
			"lastName":  lastName,
			"email":     email,
		}

		data["errorMessage"] = "A member with this email address already exists."
		render()
		return
	}

	// The passwords don't match
	if password != reenterPassword {
		data["user"] = map[string]interface{}{
			"firstName": firstName,
			"lastName":  lastName,
			"email":     email,
		}

		data["errorMessage"] = "The passwords you provided don't match. Please re-type them and try submitting again."
		render()
		return
	}

	// Create the member
	member = Member{
		Approved:  false,
		AvatarURL: "",
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Password:  passwords.HashedPasswordString(password),
	}

	if err = fa.MemberService.CreateMember(&member); err != nil {
		fa.Logger.WithError(err).Error("error creating new member")
		http.Redirect(w, r, UnexpectedErrorPath, http.StatusFound)
		return
	}

	http.Redirect(w, r, SiteAuthAccountPendingPath, http.StatusFound)
}
