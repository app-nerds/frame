package framemember

import (
	"fmt"

	"github.com/app-nerds/frame/pkg/database"
	"gorm.io/gorm"
)

type MemberServiceConfig struct {
	DB       *gorm.DB
	PageSize int
}

type MemberService struct {
	db       *gorm.DB
	pageSize int
}

func NewMemberService(config MemberServiceConfig) MemberService {
	return MemberService{
		db:       config.DB,
		pageSize: config.PageSize,
	}
}

func (s MemberService) ActivateMember(id uint) error {
	member := Member{}
	queryResult := s.db.First(&member, id)

	if queryResult.Error != nil {
		return queryResult.Error
	}

	member.Status = MembersStatus{
		ID:     MemberActiveID,
		Status: MemberActive,
	}

	queryResult = s.db.Save(&member)
	return queryResult.Error
}

func (s MemberService) CreateMember(member *Member) error {
	member.Password = member.Password.Hash()
	dbResult := s.db.Create(&member)
	return dbResult.Error
}

func (s MemberService) GetMemberByEmail(email string, includeDeleted bool) (Member, error) {
	result := Member{}

	query := s.db

	if includeDeleted {
		query = query.Unscoped()
	}

	queryResult := query.Joins("Status").Where("email = ?", email).First(&result)
	return result, queryResult.Error
}

func (s MemberService) GetMemberByEmailAndExternalID(email, id string) (Member, error) {
	result := Member{}

	queryResult := s.db.Where("email = ? AND external_id = ?", email, id).First(&result)
	return result, queryResult.Error
}

func (s MemberService) GetMemberByID(id int) (Member, error) {
	result := Member{
		Model: gorm.Model{
			ID: uint(id),
		},
	}

	queryResult := s.db.Joins("Status").First(&result)
	return result, queryResult.Error
}

func (s MemberService) GetMembers(page int) ([]Member, error) {
	result := []Member{}

	queryResult := s.db.Unscoped().Scopes(database.Paginate(page, s.pageSize)).Joins("Status").Find(&result)
	return result, queryResult.Error
}

func (s MemberService) InactivateMember(id uint) error {
	member := Member{}
	queryResult := s.db.First(&member, id)

	if queryResult.Error != nil {
		return queryResult.Error
	}

	member.Status = MembersStatus{
		ID:     MemberInactiveID,
		Status: MemberInactive,
	}

	queryResult = s.db.Save(&member)
	return queryResult.Error
}

func (s MemberService) UpdateMember(member Member) error {
	var (
		queryResult *gorm.DB
	)

	existingMember := Member{}
	existingMember.ID = member.ID

	queryResult = s.db.First(&existingMember)

	if queryResult.Error != nil {
		return fmt.Errorf("error querying for existing member in UpdateMember: %w", queryResult.Error)
	}

	/*
	 * If we've got a password in the new member struct, we are changing it
	 */
	if existingMember.Password != member.Password {
		member.Password = member.Password.Hash()
	}

	queryResult = s.db.Save(&member)
	return queryResult.Error
}
