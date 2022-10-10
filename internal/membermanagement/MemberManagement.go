package membermanagement

import (
	"net/http"

	"github.com/app-nerds/frame/internal/routepaths"
	webapp "github.com/app-nerds/frame/internal/web-app"
	"github.com/app-nerds/frame/pkg/framemember"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type InternalMemberManagementConfig struct {
	AppName       string
	Logger        *logrus.Entry
	MemberService *framemember.MemberService
	WebApp        *webapp.WebApp
}

type MemberManagement struct {
	appName       string
	logger        *logrus.Entry
	memberService *framemember.MemberService
	webApp        *webapp.WebApp
}

func NewMemberManagement(internalConfig InternalMemberManagementConfig) *MemberManagement {
	result := &MemberManagement{
		appName:       internalConfig.AppName,
		logger:        internalConfig.Logger,
		memberService: internalConfig.MemberService,
		webApp:        internalConfig.WebApp,
	}

	return result
}

func (mm *MemberManagement) RegisterRoutes(router *mux.Router) {
	router.HandleFunc(routepaths.MemberSignUpPath, mm.handleMemberSignup).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc(routepaths.MemberApiCurrentMember, mm.handleMemberCurrent).Methods(http.MethodGet)
	router.HandleFunc(routepaths.MemberApiLogOut, mm.handleMemberLogout).Methods(http.MethodGet)
	router.HandleFunc("/admin/members/manage", mm.handleAdminMembersManage).Methods(http.MethodGet)
	router.HandleFunc("/admin/api/members", mm.handleAdminApiGetMembers).Methods(http.MethodGet)
	router.HandleFunc("/admin/api/member/activate", mm.handleMemberActivate).Methods(http.MethodPut)
}
