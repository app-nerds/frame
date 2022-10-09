package frame

import (
	"github.com/app-nerds/kit/v6/passwords"
	"gorm.io/gorm"
)

type MemberStatus string

const (
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
	StatusID   int                            `json:"-"`
	Status     MembersStatus                  `json:"memberStatus"`
}

type MembersStatus struct {
	ID     int          `json:"id"`
	Status MemberStatus `json:"status"`
}
