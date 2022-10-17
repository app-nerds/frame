package framemember

import (
	"github.com/app-nerds/kit/v6/passwords"
	"gorm.io/gorm"
)

type MemberStatus string

const (
	BaseMemberRole          string       = "Member"
	MemberPendingApprovalID int          = 1
	MemberPendingApproval   MemberStatus = "Pending Approval"
	MemberActiveID          int          = 2
	MemberActive            MemberStatus = "Active"
	MemberInactiveID        int          = 3
	MemberInactive          MemberStatus = "Inactive"
)

type Member struct {
	gorm.Model

	AvatarURL  string                         `json:"avatarURL"`
	Email      string                         `json:"email"`
	ExternalID string                         `json:"-"`
	FirstName  string                         `json:"firstName"`
	LastName   string                         `json:"lastName"`
	Password   passwords.HashedPasswordString `json:"-"`
	RoleID     int                            `json:"-"`
	Role       MemberRole                     `json:"role"`
	StatusID   int                            `json:"-"`
	Status     MembersStatus                  `json:"memberStatus"`
}

type MemberRole struct {
	gorm.Model

	Color    string `json:"color"`
	RoleName string `json:"roleName"`
}

type MembersStatus struct {
	ID     int          `json:"id"`
	Status MemberStatus `json:"status"`
}
