package frame

import (
	"fmt"

	"github.com/app-nerds/frame/pkg/framemember"
)

func (fa *FrameApplication) seedDataMemberStatuses() error {
	var (
		err   error
		count int64
	)

	if err = fa.DB.Model(&framemember.MembersStatus{}).Count(&count).Error; err != nil {
		return fmt.Errorf("error getting count of member status records: %w", err)
	}

	if count > 0 {
		return nil
	}

	statuses := []framemember.MembersStatus{
		{Status: framemember.MemberPendingApproval},
		{Status: framemember.MemberActive},
		{Status: framemember.MemberInactive},
	}

	err = fa.DB.Create(&statuses).Error
	return err
}
