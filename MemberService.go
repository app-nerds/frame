package frame

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackskj/carta"
)

type MemberServiceConfig struct {
	DB       *sql.DB
	PageSize int
}

type MemberService struct {
	db       *sql.DB
	pageSize int
}

func NewMemberService(config MemberServiceConfig) MemberService {
	return MemberService{
		db:       config.DB,
		pageSize: config.PageSize,
	}
}

func (s MemberService) ActivateMember(id string) error {
	query := `
		UPDATE members SET 
			status_id = $1,
			updated_at = $2
		WHERE id = $3
	`

	_, err := s.db.Exec(query, MemberActiveID, time.Now().UTC(), id)
	return err
}

func (s MemberService) CreateMember(member *Member) error {
	id := uuid.NewString()
	member.Password = member.Password.Hash()

	query := `
		INSERT INTO members (
			id,
			created_at,
			avatar_url,
			email,
			external_id,
			first_name,
			last_name,
			password,
			role_id,
			status_id
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10
		)
	`

	_, err := s.db.Exec(
		query,
		id,
		time.Now().UTC(),
		member.AvatarURL,
		member.Email,
		member.ExternalID,
		member.FirstName,
		member.LastName,
		member.Password,
		member.RoleID,
		member.StatusID,
	)

	return err
}

func (s MemberService) DeleteMember(id string) error {
	query := `
		UPDATE members SET
			deleted_at = $1
		WHERE id = $2
	`

	_, err := s.db.Exec(query, time.Now().UTC(), id)
	return err
}

func (s MemberService) GetMemberByEmail(email string, includeDeleted bool) (Member, error) {
	query := `
		SELECT
			m.*,
			ms.status, 
			mr.role_name
		FROM members m
			INNER JOIN member_statuses ms ON m.status_id = ms.id
			INNER JOIN member_roles mr ON m.role_id = mr.id
		WHERE 1=1
			AND m.email = $1
	`

	if !includeDeleted {
		query += " AND m.deleted_at IS NULL"
	}

	rows, err := s.db.Query(query, email)

	if err != nil {
		return Member{}, err
	}

	defer rows.Close()

	members := []Member{}

	if err = carta.Map(rows, &members); err != nil {
		return Member{}, err
	}

	if len(members) < 1 {
		return Member{}, fmt.Errorf("member not found")
	}

	return members[0], nil
}

func (s MemberService) GetMemberByID(id int, includeDeleted bool) (Member, error) {
	query := `
		SELECT
			m.*,
			ms.status, 
			mr.role_name
		FROM members m
			INNER JOIN member_statuses ms ON m.status_id = ms.id
			INNER JOIN member_roles mr ON m.role_id = mr.id
		WHERE 1=1
			AND m.id = $1
	`

	if !includeDeleted {
		query += " AND m.deleted_at IS NULL"
	}

	rows, err := s.db.Query(query, id)

	if err != nil {
		return Member{}, err
	}

	defer rows.Close()

	members := []Member{}

	if err = carta.Map(rows, &members); err != nil {
		return Member{}, err
	}

	if len(members) < 1 {
		return Member{}, fmt.Errorf("member not found")
	}

	return members[0], nil
}

func (s MemberService) GetMembers(page int, includeDeleted bool) ([]Member, error) {
	members := []Member{}

	query := `
		SELECT
			m.*,
			ms.status, 
			mr.role_name
		FROM members m
			INNER JOIN member_statuses ms ON m.status_id = ms.id
			INNER JOIN member_roles mr ON m.role_id = mr.id
		WHERE 1=1
	`

	if !includeDeleted {
		query += " AND m.deleted_at IS NULL"
	}

	query += GetDBPaging(page, s.pageSize)

	rows, err := s.db.Query(query)

	if err != nil {
		return members, err
	}

	defer rows.Close()

	if err = carta.Map(rows, &members); err != nil {
		return members, err
	}

	if len(members) < 1 {
		return members, fmt.Errorf("member not found")
	}

	return members, nil
}

func (s MemberService) GetMemberRole(name string) (MemberRole, error) {
	var (
		err   error
		rows  *sql.Rows
		roles []MemberRole
	)

	query := `
		SELECT 
			id,
			created_at,
			updated_at,
			deleted_at,
			color,
			role
		FROM member_role
		WHERE role = $1
	`

	if rows, err = s.db.Query(query, name); err != nil {
		return MemberRole{}, err
	}

	defer rows.Close()

	if err = carta.Map(rows, &roles); err != nil {
		return MemberRole{}, err
	}

	if len(roles) < 1 {
		return MemberRole{}, fmt.Errorf("role not found")
	}

	return roles[0], nil
}

func (s MemberService) GetMemberRoleByID(id int) (MemberRole, error) {
	var (
		err   error
		rows  *sql.Rows
		roles []MemberRole
	)

	query := `
		SELECT 
			id,
			created_at,
			updated_at,
			deleted_at,
			color,
			role
		FROM member_role
		WHERE id = $1
	`

	if rows, err = s.db.Query(query, id); err != nil {
		return MemberRole{}, err
	}

	defer rows.Close()

	if err = carta.Map(rows, &roles); err != nil {
		return MemberRole{}, err
	}

	if len(roles) < 1 {
		return MemberRole{}, fmt.Errorf("role not found")
	}

	return roles[0], nil
}

func (s MemberService) GetMemberRoles() ([]MemberRole, error) {
	var (
		err   error
		rows  *sql.Rows
		roles []MemberRole
	)

	query := `
		SELECT 
			id,
			created_at,
			updated_at,
			deleted_at,
			color,
			role
		FROM member_role
		WHERE 1=1
	`

	if rows, err = s.db.Query(query); err != nil {
		return []MemberRole{}, err
	}

	defer rows.Close()

	if err = carta.Map(rows, &roles); err != nil {
		return []MemberRole{}, err
	}

	return roles, nil
}

func (s MemberService) CreateMemberRole(role MemberRole) (MemberRole, error) {
	var (
		err       error
		sqlResult sql.Result
		newID     int64
	)

	query := `
		INSERT INTO member_role (
			created_at,
			color,
			role
		) VALUES (
			$1,
			$2,
			$3
		)
	`

	if sqlResult, err = s.db.Exec(query, time.Now().UTC(), role.Color, role.Role); err != nil {
		return MemberRole{}, err
	}

	if newID, err = sqlResult.LastInsertId(); err != nil {
		return MemberRole{}, err
	}

	role.ID = uint(newID)
	return role, nil
}

func (s MemberService) InactivateMember(id uint) error {
	var (
		err error
	)

	query := `
		UPDATE member SET 
			status_id = $1
		WHERE id = $2
	`

	if _, err = s.db.Exec(query, MemberInactiveID, id); err != nil {
		return err
	}

	return nil
}

func (s MemberService) UpdateMember(member Member) error {
	var (
		err error
	)

	query := `
		UPDATE member SET
			update_at = $1,
			avatar_url = $2,
			email = $3,
			external_id = $4,
			first_name = $5,
			last_name = $6,
			role_id = $7,
			status_id = $8
	`

	if member.Password != "" {
		query += ", password = $9"
	}

	query += " WHERE id = $10"

	/*
	 * If we've got a password in the new member struct, we are changing it
	 */
	if member.Password != "" {
		member.Password = member.Password.Hash()
	}

	params := []interface{}{
		time.Now().UTC(),
		member.AvatarURL,
		member.Email,
		member.ExternalID,
		member.FirstName,
		member.LastName,
		member.RoleID,
		member.StatusID,
	}

	if member.Password != "" {
		params = append(params, member.Password)
	}

	params = append(params, member.ID)

	if _, err = s.db.Exec(query, params...); err != nil {
		return err
	}

	return nil
}
