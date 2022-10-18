package membermanagement

import (
	"github.com/app-nerds/frame/internal/baseviewmodel"
	"github.com/app-nerds/frame/pkg/framemember"
)

type MembersManageData struct {
	baseviewmodel.BaseViewModel
}

type MembersEditData struct {
	baseviewmodel.BaseViewModel
	Member  framemember.Member
	Message string
	Success bool
}
type MemberProfileData struct {
	baseviewmodel.BaseViewModel
	EditAvatarPath string
	Member         framemember.Member
	Message        string
	Success        bool
}

type EditAvatarData struct {
	baseviewmodel.BaseViewModel
	Member  framemember.Member
	Message string
	Success bool
}

type RolesManageData struct {
	baseviewmodel.BaseViewModel
	Roles []framemember.MemberRole
}

type RolesCreateData struct {
	baseviewmodel.BaseViewModel
	Role    framemember.MemberRole
	Success bool
	Message string
}

type RolesEditData struct {
	baseviewmodel.BaseViewModel
	Role    framemember.MemberRole
	Success bool
	Message string
}
