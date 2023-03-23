package frame

import (
	"time"

	"github.com/app-nerds/kit/v6/passwords"
)

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
	ID         string                         `json:"id"`
	CreatedAt  time.Time                      `json:"createdAt"`
	UpdatedAt  *time.Time                     `json:"updatedAt"`
	DeletedAt  *time.Time                     `json:"deletedAt"`
	AvatarURL  string                         `json:"avatarURL"`
	Email      string                         `json:"email"`
	ExternalID string                         `json:"-"`
	FirstName  string                         `json:"firstName"`
	LastName   string                         `json:"lastName"`
	Password   passwords.HashedPasswordString `json:"-"`
	RoleID     uint                           `json:"-"`
	Role       MemberRole                     `json:"role"`
	StatusID   uint                           `json:"-"`
	Status     MembersStatus                  `json:"memberStatus"`
}

type MemberRole struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
	Color     string     `json:"color"`
	Role      string     `json:"role"`
}

type MembersStatus struct {
	ID     uint         `json:"id"`
	Status MemberStatus `json:"status"`
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
