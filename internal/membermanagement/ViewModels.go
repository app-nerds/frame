package membermanagement

import (
	"github.com/app-nerds/frame/internal/baseviewmodel"
	"github.com/app-nerds/frame/pkg/framemember"
)

type MembersManageData struct {
	baseviewmodel.BaseViewModel
}

type MemberProfileData struct {
	baseviewmodel.BaseViewModel
	Member  framemember.Member
	Message string
	Success bool
}
