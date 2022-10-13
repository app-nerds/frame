package membermanagement

import (
	"net/http"

	"github.com/app-nerds/frame/internal/routepaths"
	webapp "github.com/app-nerds/frame/internal/web-app"
	"github.com/app-nerds/frame/pkg/framemember"
	"github.com/app-nerds/gobucket/v2/cmd/gobucketgo"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type InternalMemberManagementConfig struct {
	AppName        string
	GobucketClient *gobucketgo.GoBucket
	Logger         *logrus.Entry
	MemberService  *framemember.MemberService
	WebApp         *webapp.WebApp
}

type MemberManagement struct {
	appName        string
	gobucketClient *gobucketgo.GoBucket
	logger         *logrus.Entry
	memberService  *framemember.MemberService
	webApp         *webapp.WebApp
}

func NewMemberManagement(internalConfig InternalMemberManagementConfig) *MemberManagement {
	result := &MemberManagement{
		appName:        internalConfig.AppName,
		gobucketClient: internalConfig.GobucketClient,
		logger:         internalConfig.Logger,
		memberService:  internalConfig.MemberService,
		webApp:         internalConfig.WebApp,
	}

	return result
}

func (mm *MemberManagement) RegisterRoutes(router *mux.Router, adminRouter *mux.Router) {
	router.HandleFunc(routepaths.MemberSignUpPath, mm.handleMemberSignup).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc(routepaths.MemberApiCurrentMember, mm.handleMemberCurrent).Methods(http.MethodGet)
	router.HandleFunc(routepaths.MemberApiLogOut, mm.handleMemberLogout).Methods(http.MethodGet)
	router.HandleFunc(routepaths.MemberProfilePath, mm.handleMemberProfile).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc(routepaths.MemberProfileAvatarPath, mm.handleEditAvatar).Methods(http.MethodGet, http.MethodPost)
	adminRouter.HandleFunc("/members/manage", mm.handleAdminMembersManage).Methods(http.MethodGet)
	adminRouter.HandleFunc("/api/members", mm.handleAdminApiGetMembers).Methods(http.MethodGet)
	adminRouter.HandleFunc("/api/member/activate", mm.handleMemberActivate).Methods(http.MethodPut)
	adminRouter.HandleFunc("/api/member/delete/{id}", mm.handleMemberDelete).Methods(http.MethodDelete)
}
