package membermanagement

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/app-nerds/frame/internal/baseviewmodel"
	"github.com/app-nerds/frame/internal/routepaths"
	"github.com/app-nerds/frame/pkg/framemember"
	"github.com/app-nerds/frame/pkg/httputils"
	"github.com/app-nerds/frame/pkg/paging"
	webapp "github.com/app-nerds/frame/pkg/web-app"
	"github.com/app-nerds/gobucket/v2/pkg/requestcontracts"
	"github.com/app-nerds/gobucket/v2/pkg/responsecontracts"
	"github.com/app-nerds/kit/v6/passwords"
	"github.com/gorilla/mux"
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
		EditAvatarPath: routepaths.MemberProfileAvatarPath,
		Member:         framemember.Member{},
		Message:        "",
		Success:        true,
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

func (mm *MemberManagement) handleEditAvatar(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		file     multipart.File
		header   *multipart.FileHeader
		imageURL string
	)

	ctx := r.Context()
	memberEmail, _ := ctx.Value("email").(string)
	memberFirstName, _ := ctx.Value("firstName").(string)
	memberLastName, _ := ctx.Value("lastName").(string)

	data := EditAvatarData{
		BaseViewModel: baseviewmodel.BaseViewModel{
			JavascriptIncludes: webapp.JavascriptIncludes{
				{Type: "module", Src: "/frame-static/js/member-edit-avatar.js"},
			},
			AppName: mm.appName,
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

	/*
	 * Handle form post
	 */
	if r.Method == http.MethodPost {
		// Allow a 1MB post
		if err = r.ParseMultipartForm(1 << 20); err != nil {
			data.Success = false
			data.Message = "Please choose an image no larger that 250KB"
			goto rendereditavatar
		}

		// If we have no error, keep going
		if file, header, err = r.FormFile("imageFile"); err != nil {
			data.Success = false
			data.Message = "There was an error getting the file information."
			goto rendereditavatar
		}

		defer file.Close()

		/*
		 * I am supporting two types of uploads: file system, and GoBucket. If GoBucket is
		 * configured use it.
		 */
		if mm.gobucketClient != nil {
			var createImageResponse *responsecontracts.CreateImageResponse

			createImageRequest := &requestcontracts.CreateImageRequest{
				Author:       fmt.Sprintf("%s %s", memberFirstName, memberLastName),
				Bucket:       "avatars",
				Caption:      fmt.Sprintf("Avatar for %s %s", memberFirstName, memberLastName),
				DateTime:     time.Now().Format("2006-01-02T15:03:04"),
				FileContents: file,
				FileName:     header.Filename,
				Metadata:     map[string]string{},
				Name:         fmt.Sprintf("avatar-%s-%s", memberFirstName, memberLastName),
				ScaleImage:   false,
				SortIndex:    0,
				Tags:         []string{},
			}

			if createImageResponse, err = mm.gobucketClient.CreateImage(createImageRequest); err != nil {
				mm.logger.WithError(err).Error("error uploading image to Gobucket")

				data.Success = false
				data.Message = "There was an error uploading your image. Please try again."
				goto rendereditavatar
			}

			imageURL = createImageResponse.UploadedImages[0].URL
		}

		// TODO: Add local file upload support

		// Update member record
		data.Member.AvatarURL = imageURL

		if err = mm.memberService.UpdateMember(data.Member); err != nil {
			mm.logger.WithError(err).Error("error updating member after image upload")

			data.Success = false
			data.Message = "There was a problem updating your member record. Please try again."
			goto rendereditavatar
		}

		data.Success = true
		data.Message = "Avatar uploaded successfully!"
	}

rendereditavatar:
	mm.webApp.RenderTemplate(w, "member-edit-avatar.tmpl", data)
}

func (mm *MemberManagement) handleAdminApiGetMembers(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		members []framemember.Member
	)

	page := paging.GetPageFromRequest(r)

	if members, err = mm.memberService.GetMembers(page, false); err != nil {
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
	var (
		err    error
		member framemember.Member
	)

	ctx := r.Context()
	email := ctx.Value("email").(string)

	if member, err = mm.memberService.GetMemberByEmail(email, false); err != nil {
		mm.logger.WithError(err).Error("error getting member in handleMemberCurrent()")
		httputils.WriteJSON(w, http.StatusInternalServerError, httputils.CreateGenericErrorResponse("Error retrieving member information", err.Error(), ""))
		return
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

func (mm *MemberManagement) handleMemberDelete(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		id     int
		member framemember.Member
	)

	vars := mux.Vars(r)
	idString := vars["id"]

	if id, err = strconv.Atoi(idString); err != nil {
		mm.logger.WithError(err).Error("invalid member ID")
		httputils.WriteJSON(w, http.StatusBadRequest, httputils.CreateGenericErrorResponse("Invalid member ID", "", ""))
		return
	}

	if member, err = mm.memberService.GetMemberByID(id); err != nil {
		mm.logger.WithError(err).Error("error getting member in handleMemberDelete()")
		httputils.WriteJSON(w, http.StatusInternalServerError, httputils.CreateGenericErrorResponse("Error retrieving member information", err.Error(), ""))
		return
	}

	if err = mm.memberService.DeleteMember(member); err != nil {
		mm.logger.WithError(err).WithField("memberID", id).Error("error deleting member")
		httputils.WriteJSON(w, http.StatusInternalServerError, httputils.CreateGenericErrorResponse("Error deleting member", err.Error(), ""))
		return
	}

	httputils.WriteJSON(w, http.StatusOK, httputils.CreateGenericSuccessResponse("Member deleted successfully"))
}
