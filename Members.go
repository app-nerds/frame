package frame

import (
	"database/sql"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/app-nerds/gobucket/v2/cmd/gobucketgo"
	"github.com/app-nerds/gobucket/v2/pkg/requestcontracts"
	"github.com/app-nerds/gobucket/v2/pkg/responsecontracts"
	"github.com/app-nerds/kit/v6/passwords"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jackskj/carta"
	"github.com/sirupsen/logrus"
)

/*******************************************************************************
 * Internal Member Management
 ******************************************************************************/

type InternalMemberManagementConfig struct {
	AppName                  string
	CustomMemberSignupConfig *CustomMemberSignupConfig
	GobucketClient           *gobucketgo.GoBucket
	Logger                   *logrus.Entry
	MemberService            *MemberService
	WebApp                   *WebApp
}

type MemberManagement struct {
	appName                  string
	customMemberSignupConfig *CustomMemberSignupConfig
	gobucketClient           *gobucketgo.GoBucket
	logger                   *logrus.Entry
	memberService            *MemberService
	webApp                   *WebApp
}

func NewMemberManagement(internalConfig InternalMemberManagementConfig) *MemberManagement {
	result := &MemberManagement{
		appName:                  internalConfig.AppName,
		customMemberSignupConfig: internalConfig.CustomMemberSignupConfig,
		gobucketClient:           internalConfig.GobucketClient,
		logger:                   internalConfig.Logger,
		memberService:            internalConfig.MemberService,
		webApp:                   internalConfig.WebApp,
	}

	return result
}

/*******************************************************************************
 * Registration functions
 ******************************************************************************/

func (mm *MemberManagement) RegisterRoutes(router *mux.Router, adminRouter *mux.Router) {
	if mm.customMemberSignupConfig != nil {
		router.HandleFunc(MemberSignUpPath, mm.customMemberSignupConfig.Handler).Methods(http.MethodGet, http.MethodPost)
	} else {
		router.HandleFunc(MemberSignUpPath, mm.handleMemberSignup).Methods(http.MethodGet, http.MethodPost)
	}

	router.HandleFunc(MemberApiCurrentMember, mm.handleMemberCurrent).Methods(http.MethodGet)
	router.HandleFunc(MemberApiLogOut, mm.handleMemberLogout).Methods(http.MethodGet)
	router.HandleFunc(MemberProfilePath, mm.handleMemberProfile).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc(MemberProfileAvatarPath, mm.handleEditAvatar).Methods(http.MethodGet, http.MethodPost)
	adminRouter.HandleFunc("/members/manage", mm.handleAdminMembersManage).Methods(http.MethodGet)
	adminRouter.HandleFunc("/members/edit/{id}", mm.handleAdminMembersEdit).Methods(http.MethodGet, http.MethodPost)
	adminRouter.HandleFunc("/roles/manage", mm.handleAdminRolesManage).Methods(http.MethodGet)
	adminRouter.HandleFunc("/roles/create", mm.handleAdminRolesCreate).Methods(http.MethodGet, http.MethodPost)
	adminRouter.HandleFunc("/roles/edit/{id}", mm.handleAdminRolesEdit).Methods(http.MethodGet, http.MethodPost)
	adminRouter.HandleFunc("/api/members", mm.handleAdminApiGetMembers).Methods(http.MethodGet)
	adminRouter.HandleFunc("/api/member/activate", mm.handleMemberActivate).Methods(http.MethodPut)
	adminRouter.HandleFunc("/api/member/delete/{id}", mm.handleMemberDelete).Methods(http.MethodDelete)
	adminRouter.HandleFunc("/api/member/role", mm.handleGetMemberRoles).Methods(http.MethodGet)
}

func (mm *MemberManagement) RegisterAdminTemplate() TemplateCollection {
	result := TemplateCollection{}

	result = append(result, Template{Name: "admin-members-manage.tmpl", IsLayout: false, UseLayout: "admin-layout.tmpl"})
	result = append(result, Template{Name: "admin-members-edit.tmpl", IsLayout: false, UseLayout: "admin-layout.tmpl"})
	result = append(result, Template{Name: "admin-roles-manage.tmpl", IsLayout: false, UseLayout: "admin-layout.tmpl"})
	result = append(result, Template{Name: "admin-roles-create.tmpl", IsLayout: false, UseLayout: "admin-layout.tmpl"})
	result = append(result, Template{Name: "admin-roles-edit.tmpl", IsLayout: false, UseLayout: "admin-layout.tmpl"})

	return result
}

func (mm *MemberManagement) RegisterTemplates() TemplateCollection {
	result := TemplateCollection{}

	result = append(result, Template{Name: "member-profile.tmpl", IsLayout: false, UseLayout: "layout.tmpl"})
	result = append(result, Template{Name: "member-edit-avatar.tmpl", IsLayout: false, UseLayout: "layout.tmpl"})

	return result
}

/*******************************************************************************
 * Handlers
 ******************************************************************************/

func (mm *MemberManagement) handleAdminMembersManage(w http.ResponseWriter, r *http.Request) {
	data := MembersManageData{
		BaseViewModel: BaseViewModel{
			JavascriptIncludes: JavascriptIncludes{
				{Type: "module", Src: "/pages/admin-members-manage.js"},
			},
			AppName: mm.appName,
		},
	}

	mm.webApp.RenderTemplate(w, "admin-members-manage.tmpl", data)
}

func (mm *MemberManagement) handleAdminMembersEdit(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		id  string
	)

	vars := mux.Vars(r)
	id = vars["id"]

	data := MembersEditData{
		BaseViewModel: BaseViewModel{
			JavascriptIncludes: JavascriptIncludes{
				{Type: "module", Src: "/pages/admin-members-edit.js"},
			},
			AppName: mm.appName,
		},
	}

	if data.Member, err = mm.memberService.GetMemberByID(id, false); err != nil {
		mm.logger.WithError(err).Error("error retrieving member in handleAdminMembersEdit")

		data.Success = false
		data.Message = "There was a problem retrieving this member's information. Please try again."
		goto rendermembersedit
	}

	/*
	 * POST
	 */
	if r.Method == http.MethodPost {
		data.Member.Password = ""
		_ = r.ParseForm()

		if data.Member.AvatarURL == "" {
			data.Member.AvatarURL = "/frame-static/images/blank-profile-picture.png"
		}

		if r.FormValue("lastName") == "" {
			data.Success = false
			data.Message = "Please provide a last name."
			goto rendermembersedit
		}

		if r.FormValue("firstName") == "" {
			data.Success = false
			data.Message = "Please provide a first name."
			goto rendermembersedit
		}

		roleID, err := strconv.Atoi(r.FormValue("role"))

		if err != nil {
			data.Success = false
			data.Message = "Invalid Role selected"
			goto rendermembersedit
		}

		role, err := mm.memberService.GetMemberRoleByID(roleID)

		if err != nil {
			data.Success = false
			data.Message = "There was a problem getting role information"
			goto rendermembersedit
		}

		data.Member.FirstName = r.FormValue("firstName")
		data.Member.LastName = r.FormValue("lastName")
		data.Member.Role = role

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

rendermembersedit:
	mm.webApp.RenderTemplate(w, "admin-members-edit.tmpl", data)
}

func (mm *MemberManagement) handleMemberProfile(w http.ResponseWriter, r *http.Request) {
	var (
		err error
	)

	ctx := r.Context()
	memberEmail, _ := ctx.Value("email").(string)

	data := MemberProfileData{
		BaseViewModel: BaseViewModel{
			JavascriptIncludes: JavascriptIncludes{},
			AppName:            mm.appName,
			Stylesheets: []string{
				"/frame-static/css/frame-page-styles.css",
			},
		},
		EditAvatarPath: MemberProfileAvatarPath,
		Member:         Member{},
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
		BaseViewModel: BaseViewModel{
			JavascriptIncludes: JavascriptIncludes{
				{Type: "module", Src: "/frame-static/js/member-edit-avatar.js"},
			},
			AppName: mm.appName,
			Stylesheets: []string{
				"/frame-static/css/frame-page-styles.css",
			},
		},
		Member:  Member{},
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
		members []Member
	)

	page := GetPageFromRequest(r)

	if members, err = mm.memberService.GetMembers(page, false); err != nil {
		mm.logger.WithError(err).Error("error getting members")
		WriteJSON(w, http.StatusInternalServerError, CreateGenericErrorResponse("There was a problem retrieving members", err.Error(), ""))
		return
	}

	WriteJSON(w, http.StatusOK, members)
}

/*
PUT /admin/api/member/activate
*/
func (mm *MemberManagement) handleMemberActivate(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		id  string
	)

	r.ParseForm()

	id = r.FormValue("id")

	if err = mm.memberService.ActivateMember(id); err != nil {
		mm.logger.WithError(err).Error("error activating member")
		WriteJSON(w, http.StatusInternalServerError, CreateGenericErrorResponse("Error activating member", err.Error(), ""))
		return
	}

	WriteJSON(w, http.StatusOK, CreateGenericSuccessResponse("Member activated!"))
}

/*
GET /api/member/current
*/
func (mm *MemberManagement) handleMemberCurrent(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		member Member
	)

	ctx := r.Context()
	email := ctx.Value("email").(string)

	if member, err = mm.memberService.GetMemberByEmail(email, false); err != nil {
		mm.logger.WithError(err).Error("error getting member in handleMemberCurrent()")
		WriteJSON(w, http.StatusInternalServerError, CreateGenericErrorResponse("Error retrieving member information", err.Error(), ""))
		return
	}

	WriteJSON(w, http.StatusOK, member)
}

/*
GET /api/member/logout
*/
func (mm *MemberManagement) handleMemberLogout(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		session *sessions.Session
	)

	if session, err = mm.webApp.GetSessionStore().Get(r, mm.webApp.GetSessionName()); err != nil {
		mm.logger.WithError(err).Error("error getting session information")
		WriteJSON(w, http.StatusInternalServerError, CreateGenericErrorResponse("error getting session information", err.Error(), ""))
		return
	}

	session.Options.MaxAge = -1

	if err = mm.webApp.GetSessionStore().Save(r, w, session); err != nil {
		mm.logger.WithError(err).Error("error deleting session")
		WriteJSON(w, http.StatusInternalServerError, CreateGenericErrorResponse("error deleting session", err.Error(), ""))
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
		member Member
		role   MemberRole
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

	// Get the base member role
	if role, err = mm.memberService.GetMemberRole(BaseMemberRole); err != nil {
		mm.logger.WithError(err).Error("error retrieving member role in handleMemberSignup()")

		data.User.FirstName = firstName
		data.User.LastName = lastName
		data.User.Email = email
		data.ErrorMessage = "There was a problem getting some information before creating your member. Please try again."
		render()
		return
	}

	// Create the member
	member = Member{
		AvatarURL: "",
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Password:  passwords.HashedPasswordString(password),
		Status: MembersStatus{
			ID:     MemberPendingApprovalID,
			Status: MemberPendingApproval,
		},
		Role: role,
	}

	if err = mm.memberService.CreateMember(&member); err != nil {
		mm.logger.WithError(err).Error("error creating new member")
		http.Redirect(w, r, UnexpectedErrorPath, http.StatusFound)
		return
	}

	http.Redirect(w, r, SiteAuthAccountPendingPath, http.StatusFound)
}

func (mm *MemberManagement) handleMemberDelete(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		id  string
	)

	vars := mux.Vars(r)
	id = vars["id"]

	if err = mm.memberService.DeleteMember(id); err != nil {
		mm.logger.WithError(err).WithField("memberID", id).Error("error deleting member")
		WriteJSON(w, http.StatusInternalServerError, CreateGenericErrorResponse("Error deleting member", err.Error(), ""))
		return
	}

	WriteJSON(w, http.StatusOK, CreateGenericSuccessResponse("Member deleted successfully"))
}

func (mm *MemberManagement) handleGetMemberRoles(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		roles []MemberRole
	)

	if roles, err = mm.memberService.GetMemberRoles(); err != nil {
		mm.logger.WithError(err).Error("error retrieving member roles")
		WriteJSON(w, http.StatusInternalServerError, CreateGenericErrorResponse("Error retrieving roles", "", ""))
		return
	}

	WriteJSON(w, http.StatusOK, roles)
}

func (mm *MemberManagement) handleAdminRolesManage(w http.ResponseWriter, r *http.Request) {
	var (
		err error
	)

	data := RolesManageData{
		BaseViewModel: BaseViewModel{
			JavascriptIncludes: JavascriptIncludes{
				{Type: "module", Src: "/pages/admin-roles-manage.js"},
			},
			AppName:     mm.appName,
			Stylesheets: []string{},
		},
		Roles: []MemberRole{},
	}

	if data.Roles, err = mm.memberService.GetMemberRoles(); err != nil {
		mm.logger.WithError(err).Error("error retrieving member roles in handleAdminRolesManage()")
		http.Redirect(w, r, UnexpectedErrorPath, http.StatusFound)
		return
	}

	mm.webApp.RenderTemplate(w, "admin-roles-manage.tmpl", data)
}

func (mm *MemberManagement) handleAdminRolesCreate(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		existing MemberRole
	)

	data := RolesCreateData{
		BaseViewModel: BaseViewModel{
			JavascriptIncludes: JavascriptIncludes{
				{Type: "module", Src: "/pages/admin-roles-create.js"},
			},
			AppName:     mm.appName,
			Stylesheets: []string{},
		},
		Role:    MemberRole{},
		Success: true,
		Message: "",
	}

	/*
	 * POST
	 */
	if r.Method == http.MethodPost {
		_ = r.ParseForm()

		roleName := r.FormValue("roleName")
		color := r.FormValue("color")

		if roleName == "" {
			data.Success = false
			data.Message = "Please provide a name for your new role."
			goto renderrolescreate
		}

		if color == "" {
			data.Success = false
			data.Message = "Please select a color to represent this role."
			goto renderrolescreate
		}

		// Make sure we don't have a role by this name already
		if existing, err = mm.memberService.GetMemberRole(roleName); err != nil && !errors.Is(err, sql.ErrNoRows) {
			mm.logger.WithError(err).Error("error checking for existing role by name in handleAdminRolesCreate")

			data.Success = false
			data.Message = "There was a problem getting role information. Please try again."
			goto renderrolescreate
		}

		if existing.ID > 0 {
			data.Success = false
			data.Message = "A role with this name already exists. Please choose another."
			goto renderrolescreate
		}

		data.Role.Role = roleName
		data.Role.Color = color

		if data.Role, err = mm.memberService.CreateMemberRole(data.Role); err != nil {
			mm.logger.WithError(err).Error("error creating new role in handleAdminRolesCreate")

			data.Success = false
			data.Message = "There was a problem creating your new role. Please try again."
			goto renderrolescreate
		}

		http.Redirect(w, r, "/admin/roles/manage", http.StatusFound)
		return
	}

renderrolescreate:
	mm.webApp.RenderTemplate(w, "admin-roles-create.tmpl", data)
}

func (mm *MemberManagement) handleAdminRolesEdit(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		idString string
		id       int
	)

	vars := mux.Vars(r)
	idString = vars["id"]

	data := RolesEditData{
		BaseViewModel: BaseViewModel{
			JavascriptIncludes: JavascriptIncludes{
				{Type: "module", Src: "/pages/admin-roles-edit.js"},
			},
			AppName:     mm.appName,
			Stylesheets: []string{},
		},
		Role:    MemberRole{},
		Success: true,
		Message: "",
	}

	if id, err = strconv.Atoi(idString); err != nil {
		data.Success = false
		data.Message = "Invalid role ID"
		goto renderrolesedit
	}

	/*
	 * You can't edit the first role, Member. That is system generated
	 */
	if id == 1 {
		data.Success = false
		data.Message = "You are not allowed to edit this role"
		goto renderrolesedit
	}

	data.Role, err = mm.memberService.GetMemberRoleByID(id)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		mm.logger.WithError(err).Error("error retrieving role in handleAdminRolesEdit")

		data.Success = false
		data.Message = "There was a problem retrieving role information. Please try again."
		goto renderrolesedit
	}

	/*
	 * POST
	 */
	if r.Method == http.MethodPost {
	}

renderrolesedit:
	mm.webApp.RenderTemplate(w, "admin-roles-edit.tmpl", data)
}

/*******************************************************************************
 * Services
 ******************************************************************************/

type MemberServiceConfig struct {
	DB       *sql.DB
	PageSize int
}

type MemberService struct {
	db       *sql.DB
	pageSize int
}

func NewMemberService(config MemberServiceConfig) MemberService {
	return MemberService{
		db:       config.DB,
		pageSize: config.PageSize,
	}
}

func (s MemberService) ActivateMember(id string) error {
	query := `
		UPDATE members SET 
			status_id = $1,
			updated_at = $2
		WHERE id = $3
	`

	_, err := s.db.Exec(query, MemberActiveID, time.Now().UTC(), id)
	return err
}

func (s MemberService) CreateMember(member *Member) error {
	id := uuid.NewString()
	member.Password = member.Password.Hash()

	query := `
		INSERT INTO members (
			id,
			created_at,
			avatar_url,
			email,
			external_id,
			first_name,
			last_name,
			password,
			role_id,
			status_id
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10
		)
	`

	_, err := s.db.Exec(
		query,
		id,
		time.Now().UTC(),
		member.AvatarURL,
		member.Email,
		member.ExternalID,
		member.FirstName,
		member.LastName,
		member.Password,
		member.Role.ID,
		member.Status.ID,
	)

	return err
}

func (s MemberService) DeleteMember(id string) error {
	query := `
		UPDATE members SET
			deleted_at = $1
		WHERE id = $2
	`

	_, err := s.db.Exec(query, time.Now().UTC(), id)
	return err
}

func (s MemberService) GetMemberByEmail(email string, includeDeleted bool) (Member, error) {
	query := `
		SELECT
			members.id AS member_id,
			members.created_at AS member_created_at,
			members.updated_at AS member_updated_at,
			members.deleted_at AS member_deleted_at,
			members.avatar_url AS member_avatar_url,
			members.email AS member_email,
			members.external_id AS member_external_id,
			members.first_name AS member_first_name,
			members.last_name AS member_last_name,
			members.password AS member_password,
			member_statuses.id AS status_id,
			member_statuses.status AS status_status, 
			member_roles.id AS role_id,
			member_roles.role AS role_role,
			member_roles.color
		FROM members 
			INNER JOIN member_statuses ON members.status_id = member_statuses.id
			INNER JOIN member_roles ON members.role_id = member_roles.id
		WHERE 1=1
			AND members.email = $1
	`

	if !includeDeleted {
		query += " AND members.deleted_at IS NULL"
	}

	rows, err := s.db.Query(query, email)

	if err != nil {
		return Member{}, err
	}

	defer rows.Close()

	members := []Member{}

	if err = carta.Map(rows, &members); err != nil {
		return Member{}, err
	}

	if len(members) < 1 {
		return Member{}, fmt.Errorf("member not found: %w", sql.ErrNoRows)
	}

	return members[0], nil
}

func (s MemberService) GetMemberByID(id string, includeDeleted bool) (Member, error) {
	query := `
		SELECT
			members.id AS member_id,
			members.created_at AS member_created_at,
			members.updated_at AS member_updated_at,
			members.deleted_at AS member_deleted_at,
			members.avatar_url AS member_avatar_url,
			members.email AS member_email,
			members.external_id AS member_external_id,
			members.first_name AS member_first_name,
			members.last_name AS member_last_name,
			members.password AS member_password,
			member_statuses.id AS status_id,
			member_statuses.status AS status_status, 
			member_roles.id AS role_id,
			member_roles.role AS role_role,
			member_roles.color
		FROM members 
			INNER JOIN member_statuses ON members.status_id = member_statuses.id
			INNER JOIN member_roles ON members.role_id = member_roles.id
		WHERE 1=1
			AND members.id = $1
	`

	if !includeDeleted {
		query += " AND members.deleted_at IS NULL"
	}

	rows, err := s.db.Query(query, id)

	if err != nil {
		return Member{}, err
	}

	defer rows.Close()

	members := []Member{}

	if err = carta.Map(rows, &members); err != nil {
		return Member{}, err
	}

	if len(members) < 1 {
		return Member{}, fmt.Errorf("member not found: %w", sql.ErrNoRows)
	}

	return members[0], nil
}

func (s MemberService) GetMembers(page int, includeDeleted bool) ([]Member, error) {
	members := []Member{}

	query := `
		SELECT
			members.id AS member_id,
			members.created_at AS member_created_at,
			members.updated_at AS member_updated_at,
			members.deleted_at AS member_deleted_at,
			members.avatar_url AS member_avatar_url,
			members.email AS member_email,
			members.external_id AS member_external_id,
			members.first_name AS member_first_name,
			members.last_name AS member_last_name,
			members.password AS member_password,
			member_statuses.id AS status_id,
			member_statuses.status AS status_status, 
			member_roles.id AS role_id,
			member_roles.role AS role_role,
			member_roles.color
		FROM members 
			INNER JOIN member_statuses ON members.status_id = member_statuses.id
			INNER JOIN member_roles ON members.role_id = member_roles.id
		WHERE 1=1
	`

	if !includeDeleted {
		query += " AND members.deleted_at IS NULL"
	}

	query += GetDBPaging(page, s.pageSize)

	rows, err := s.db.Query(query)

	if err != nil {
		return members, err
	}

	defer rows.Close()

	if err = carta.Map(rows, &members); err != nil {
		return members, err
	}

	return members, nil
}

func (s MemberService) GetMemberRole(name string) (MemberRole, error) {
	var (
		err   error
		rows  *sql.Rows
		roles []MemberRole
	)

	query := `
		SELECT 
			id AS role_id,
			created_at AS member_roles_created_at,
			updated_at AS member_roles_updated_at,
			deleted_at AS member_roles_deleted_at,
			color AS color,
			role AS role_role
		FROM member_roles
		WHERE role = $1
	`

	if rows, err = s.db.Query(query, name); err != nil {
		return MemberRole{}, err
	}

	defer rows.Close()

	if err = carta.Map(rows, &roles); err != nil {
		return MemberRole{}, err
	}

	if len(roles) < 1 {
		return MemberRole{}, fmt.Errorf("role not found: %w", sql.ErrNoRows)
	}

	return roles[0], nil
}

func (s MemberService) GetMemberRoleByID(id int) (MemberRole, error) {
	var (
		err   error
		rows  *sql.Rows
		roles []MemberRole
	)

	query := `
		SELECT 
			id AS role_id,
			created_at AS member_roles_created_at,
			updated_at AS member_roles_updated_at,
			deleted_at AS member_roles_deleted_at,
			color AS color,
			role AS role_role
		FROM member_roles
		WHERE id = $1
	`

	if rows, err = s.db.Query(query, id); err != nil {
		return MemberRole{}, err
	}

	defer rows.Close()

	if err = carta.Map(rows, &roles); err != nil {
		return MemberRole{}, err
	}

	if len(roles) < 1 {
		return MemberRole{}, fmt.Errorf("role not found")
	}

	return roles[0], nil
}

func (s MemberService) GetMemberRoles() ([]MemberRole, error) {
	var (
		err   error
		rows  *sql.Rows
		roles []MemberRole
	)

	query := `
		SELECT 
			id AS role_id,
			created_at AS member_roles_created_at,
			updated_at AS member_roles_updated_at,
			deleted_at AS member_roles_deleted_at,
			color AS color,
			role AS role_role
		FROM member_roles
		WHERE 1=1
	`

	if rows, err = s.db.Query(query); err != nil {
		return []MemberRole{}, err
	}

	defer rows.Close()

	if err = carta.Map(rows, &roles); err != nil {
		return []MemberRole{}, err
	}

	return roles, nil
}

func (s MemberService) CreateMemberRole(role MemberRole) (MemberRole, error) {
	var (
		err       error
		sqlResult sql.Result
		newID     int64
	)

	query := `
		INSERT INTO member_roles (
			created_at,
			color,
			role
		) VALUES (
			$1,
			$2,
			$3
		)
	`

	if sqlResult, err = s.db.Exec(query, time.Now().UTC(), role.Color, role.Role); err != nil {
		return MemberRole{}, err
	}

	if newID, err = sqlResult.LastInsertId(); err != nil {
		return MemberRole{}, err
	}

	role.ID = uint(newID)
	return role, nil
}

func (s MemberService) InactivateMember(id uint) error {
	var (
		err error
	)

	query := `
		UPDATE members SET 
			status_id = $1
		WHERE id = $2
	`

	if _, err = s.db.Exec(query, MemberInactiveID, id); err != nil {
		return err
	}

	return nil
}

func (s MemberService) UpdateMember(member Member) error {
	var (
		err error
	)

	query := `
		UPDATE members SET
			updated_at = $1,
			avatar_url = $2,
			email = $3,
			external_id = $4,
			first_name = $5,
			last_name = $6,
			role_id = $7,
			status_id = $8
	`

	if member.Password != "" {
		query += ", password = $9"
		query += " WHERE id = $10"
	} else {
		query += " WHERE id = $9"
	}

	/*
	 * If we've got a password in the new member struct, we are changing it
	 */
	if member.Password != "" {
		member.Password = member.Password.Hash()
	}

	params := []interface{}{
		time.Now().UTC(),
		member.AvatarURL,
		member.Email,
		member.ExternalID,
		member.FirstName,
		member.LastName,
		member.Role.ID,
		member.Status.ID,
	}

	if member.Password != "" {
		params = append(params, member.Password)
	}

	params = append(params, member.ID)

	if _, err = s.db.Exec(query, params...); err != nil {
		return err
	}

	return nil
}

/*******************************************************************************
 * Model
 ******************************************************************************/

type MemberStatus string

const (
	BaseMemberRole          string       = "Member"
	MemberPendingApprovalID uint         = 1
	MemberPendingApproval   MemberStatus = "Pending Approval"
	MemberActiveID          uint         = 2
	MemberActive            MemberStatus = "Active"
	MemberInactiveID        uint         = 3
	MemberInactive          MemberStatus = "Inactive"
)

type Member struct {
	ID         string                         `json:"id" db:"member_id"`
	CreatedAt  time.Time                      `json:"createdAt" db:"member_created_at"`
	UpdatedAt  *time.Time                     `json:"updatedAt" db:"member_updated_at"`
	DeletedAt  *time.Time                     `json:"deletedAt" db:"member_deleted_at"`
	AvatarURL  string                         `json:"avatarURL" db:"member_avatar_url"`
	Email      string                         `json:"email" db:"member_email"`
	ExternalID string                         `json:"-" db:"member_external_id"`
	FirstName  string                         `json:"firstName" db:"member_first_name"`
	LastName   string                         `json:"lastName" db:"member_last_name"`
	Password   passwords.HashedPasswordString `json:"-" db:"member_password"`
	Role       MemberRole                     `json:"role"`
	Status     MembersStatus                  `json:"memberStatus"`
}

type MemberRole struct {
	ID        uint       `json:"id" db:"role_id"`
	CreatedAt time.Time  `json:"createdAt" db:"member_roles_created_at"`
	UpdatedAt *time.Time `json:"updatedAt" db:"member_roles_updated_at"`
	DeletedAt *time.Time `json:"deletedAt" db:"member_roles_deleted_at"`
	Color     string     `json:"color" db:"color"`
	Role      string     `json:"role" db:"role_role"`
}

type MembersStatus struct {
	ID     uint         `json:"id" db:"status_id"`
	Status MemberStatus `json:"status" db:"status_status"`
}

type MembersManageData struct {
	BaseViewModel
}

type MembersEditData struct {
	BaseViewModel
	Member  Member
	Message string
	Success bool
}
type MemberProfileData struct {
	BaseViewModel
	EditAvatarPath string
	Member         Member
	Message        string
	Success        bool
}

type EditAvatarData struct {
	BaseViewModel
	Member  Member
	Message string
	Success bool
}

type RolesManageData struct {
	BaseViewModel
	Roles []MemberRole
}

type RolesCreateData struct {
	BaseViewModel
	Role    MemberRole
	Success bool
	Message string
}

type RolesEditData struct {
	BaseViewModel
	Role    MemberRole
	Success bool
	Message string
}
