package frame

import (
	"net/http"

	"gorm.io/gorm"
)

type MemberService struct {
	frame *FrameApplication
}

func newMemberService(frame *FrameApplication) MemberService {
	return MemberService{
		frame: frame,
	}
}

func (s MemberService) ActivateMember(id uint) error {
	member := Member{}
	queryResult := s.frame.DB.First(&member, id)

	if queryResult.Error != nil {
		return queryResult.Error
	}

	member.Status = MembersStatus{
		ID:     MemberActiveID,
		Status: MemberActive,
	}

	queryResult = s.frame.DB.Save(&member)
	return queryResult.Error
}

func (s MemberService) CreateMember(member *Member) error {
	member.Password = member.Password.Hash()
	dbResult := s.frame.DB.Create(&member)
	return dbResult.Error
}

func (s MemberService) GetMemberByEmail(email string, includeDeleted bool) (Member, error) {
	result := Member{}

	query := s.frame.DB

	if includeDeleted {
		query = query.Unscoped()
	}

	queryResult := query.Joins("Status").Where("email = ?", email).First(&result)
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

	queryResult := s.frame.DB.Joins("Status").First(&result)
	return result, queryResult.Error
}

func (s MemberService) GetMembers(r *http.Request, page int) ([]Member, error) {
	result := []Member{}

	queryResult := s.frame.DB.Unscoped().Scopes(s.frame.paginate(r)).Joins("Status").Find(&result)
	return result, queryResult.Error
}

func (s MemberService) InactivateMember(id uint) error {
	member := Member{}
	queryResult := s.frame.DB.First(&member, id)

	if queryResult.Error != nil {
		return queryResult.Error
	}

	member.Status = MembersStatus{
		ID:     MemberInactiveID,
		Status: MemberInactive,
	}

	queryResult = s.frame.DB.Save(&member)
	return queryResult.Error
}
