package frame

import (
	"github.com/app-nerds/kit/v6/passwords"
	"gorm.io/gorm"
)

type Member struct {
	gorm.Model

	Approved   bool                           `json:"approved"`
	AvatarURL  string                         `json:"avatarURL"`
	Email      string                         `json:"email"`
	ExternalID string                         `json:"-"`
	FirstName  string                         `json:"firstName"`
	LastName   string                         `json:"lastName"`
	Password   passwords.HashedPasswordString `json:"password"`
}

type MemberService struct {
	frame *FrameApplication
}

func newMemberService(frame *FrameApplication) MemberService {
	return MemberService{
		frame: frame,
	}
}

func (s MemberService) CreateMember(member *Member) error {
	member.Password = member.Password.Hash()
	dbResult := s.frame.DB.Create(&member)
	return dbResult.Error
}

func (s MemberService) GetMemberByEmail(email string) (Member, error) {
	result := Member{}

	queryResult := s.frame.DB.Where("email = ?", email).First(&result)
	return result, queryResult.Error
}

func (s MemberService) GetMemberByEmailAndExternalID(email, id string) (Member, error) {
	result := Member{}

	queryResult := s.frame.DB.Where("email = ? AND external_id = ?", email, id).First(&result)
	return result, queryResult.Error
}

func (s MemberService) GetMemberByID(id int) (Member, error) {
	result := Member{
		Model: gorm.Model{
			ID: uint(id),
		},
	}

	queryResult := s.frame.DB.First(&result)
	return result, queryResult.Error
}
