package membermanagement

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/app-nerds/frame/internal/baseviewmodel"
	"github.com/app-nerds/frame/internal/routepaths"
	"github.com/app-nerds/frame/pkg/framemember"
	"github.com/app-nerds/frame/pkg/httputils"
	"github.com/app-nerds/frame/pkg/paging"
	webapp "github.com/app-nerds/frame/pkg/web-app"
	"github.com/app-nerds/kit/v6/passwords"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

func (mm *MemberManagement) handleAdminMembersManage(w http.ResponseWriter, r *http.Request) {
	data := MembersManageData{
		BaseViewModel: baseviewmodel.BaseViewModel{
			JavascriptIncludes: webapp.JavascriptIncludes{
				{Type: "module", Src: "/pages/admin-members-manage.js"},
			},
			AppName: mm.appName,
		},
	}

	mm.webApp.RenderTemplate(w, "admin-members-manage.tmpl", data)
}

func (mm *MemberManagement) handleMemberProfile(w http.ResponseWriter, r *http.Request) {
	var (
		err error
	)

	ctx := r.Context()
	memberEmail, _ := ctx.Value("email").(string)

	data := MemberProfileData{
		BaseViewModel: baseviewmodel.BaseViewModel{
			JavascriptIncludes: webapp.JavascriptIncludes{},
			AppName:            mm.appName,
			Stylesheets: []string{
				"/frame-static/css/frame-page-styles.css",
			},
		},
		Member:  framemember.Member{},
		Message: "",
		Success: true,
	}

	if data.Member, err = mm.memberService.GetMemberByEmail(memberEmail, false); err != nil {
		mm.logger.WithError(err).Error("error getting member information in handleMemberProfile()")
		mm.webApp.UnexpectedError(w, r)
		return
	}

	if data.Member.AvatarURL == "" {
		data.Member.AvatarURL = "/frame-static/images/blank-profile-picture.png"
	}

	if r.Method == http.MethodPost {
		r.ParseForm()

		if r.FormValue("lastName") == "" {
			data.Success = false
			data.Message = "Please provide a last name."
		}

		if r.FormValue("firstName") == "" {
			data.Success = false
			data.Message = "Please provide a first name."
		}

		if data.Success {
			data.Member.FirstName = r.FormValue("firstName")
			data.Member.LastName = r.FormValue("lastName")

			if r.FormValue("password") != "" {
				data.Member.Password = passwords.HashedPasswordString(r.FormValue("password"))
			}

			if err = mm.memberService.UpdateMember(data.Member); err != nil {
				mm.logger.WithError(err).WithFields(logrus.Fields{
					"memberID": data.Member.ID,
				}).Error("error updating member")

				mm.webApp.UnexpectedError(w, r)
				return
			}

			data.Success = true
			data.Message = "Member updated successfully!"
		}
	}

	mm.webApp.RenderTemplate(w, "member-profile.tmpl", data)
}

func (mm *MemberManagement) handleAdminApiGetMembers(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		members []framemember.Member
	)

	page := paging.GetPageFromRequest(r)

	if members, err = mm.memberService.GetMembers(page); err != nil {
		mm.logger.WithError(err).Error("error getting members")
		httputils.WriteJSON(w, http.StatusInternalServerError, httputils.CreateGenericErrorResponse("There was a problem retrieving members", err.Error(), ""))
		return
	}

	httputils.WriteJSON(w, http.StatusOK, members)
}

/*
PUT /admin/api/member/activate
*/
func (mm *MemberManagement) handleMemberActivate(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		idString string
		id       int
	)

	r.ParseForm()

	idString = r.FormValue("id")

	if id, err = strconv.Atoi(idString); err != nil {
		httputils.WriteJSON(w, http.StatusBadRequest, httputils.CreateGenericErrorResponse("Invalid member ID", fmt.Sprintf("%s could not be converted to an integer", idString), ""))
		return
	}

	if err = mm.memberService.ActivateMember(uint(id)); err != nil {
		mm.logger.WithError(err).Error("error activating member")
		httputils.WriteJSON(w, http.StatusInternalServerError, httputils.CreateGenericErrorResponse("Error activating member", err.Error(), ""))
		return
	}

	httputils.WriteJSON(w, http.StatusOK, httputils.CreateGenericSuccessResponse("Member activated!"))
}

/*
GET /api/member/current
*/
func (mm *MemberManagement) handleMemberCurrent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	member := map[string]interface{}{
		"id":        ctx.Value("memberID"),
		"email":     ctx.Value("email"),
		"firstName": ctx.Value("firstName"),
		"lastName":  ctx.Value("lastName"),
		"avatarURL": ctx.Value("avatarURL"),
		"status":    ctx.Value("status"),
	}

	httputils.WriteJSON(w, http.StatusOK, member)
}

/*
/api/member/logout
*/
func (mm *MemberManagement) handleMemberLogout(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		session *sessions.Session
	)

	if session, err = mm.webApp.GetSessionStore().Get(r, mm.webApp.GetSessionName()); err != nil {
		mm.logger.WithError(err).Error("error getting session information")
		httputils.WriteJSON(w, http.StatusInternalServerError, httputils.CreateGenericErrorResponse("error getting session information", err.Error(), ""))
		return
	}

	session.Options.MaxAge = -1

	if err = mm.webApp.GetSessionStore().Save(r, w, session); err != nil {
		mm.logger.WithError(err).Error("error deleting session")
		httputils.WriteJSON(w, http.StatusInternalServerError, httputils.CreateGenericErrorResponse("error deleting session", err.Error(), ""))
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

/*
POST /member/create-account
*/
func (mm *MemberManagement) handleMemberSignup(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		member framemember.Member
	)

	data := struct {
		ErrorMessage string
		Stylesheets  []string
		User         struct {
			FirstName string
			LastName  string
			Email     string
		}
	}{
		Stylesheets: []string{
			"/frame-static/css/frame-page-styles.css",
		},
	}

	data.User.FirstName = ""
	data.User.LastName = ""
	data.User.Email = ""

	render := func() {
		mm.webApp.RenderTemplate(w, "sign-up.tmpl", data)
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

	member, err = mm.memberService.GetMemberByEmail(email, true)

	// We already have a member with this email address
	if err == nil {
		data.User.FirstName = firstName
		data.User.LastName = lastName
		data.User.Email = email
		data.ErrorMessage = "A member with this email address already exists."
		render()
		return
	}

	// The passwords don't match
	if password != reenterPassword {
		data.User.FirstName = firstName
		data.User.LastName = lastName
		data.User.Email = email
		data.ErrorMessage = "The passwords you provided don't match. Please re-type them and try submitting again."
		render()
		return
	}

	// Create the member
	member = framemember.Member{
		AvatarURL: "",
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Password:  passwords.HashedPasswordString(password),
		Status: framemember.MembersStatus{
			ID:     framemember.MemberPendingApprovalID,
			Status: framemember.MemberPendingApproval,
		},
	}

	if err = mm.memberService.CreateMember(&member); err != nil {
		mm.logger.WithError(err).Error("error creating new member")
		http.Redirect(w, r, routepaths.UnexpectedErrorPath, http.StatusFound)
		return
	}

	http.Redirect(w, r, routepaths.SiteAuthAccountPendingPath, http.StatusFound)
}
