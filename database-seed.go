package frame

import "fmt"

func (fa *FrameApplication) seedDataMemberStatuses() error {
	var (
		err   error
		count int64
	)

	if err = fa.DB.Model(&MembersStatus{}).Count(&count).Error; err != nil {
		return fmt.Errorf("error getting count of member status records: %w", err)
	}

	if count > 0 {
		return nil
	}

	statuses := []MembersStatus{
		{Status: MemberPendingApproval},
		{Status: MemberActive},
		{Status: MemberInactive},
	}

	err = fa.DB.Create(&statuses).Error
	return err
}
