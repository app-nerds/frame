package frame

import (
	"context"
	"fmt"
	"time"
)

func (f *FrameApplication) GetDBContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(f.Config.DatabaseTimeout)*time.Second)
}

func GetDBPaging(page int, pageSize int) string {
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * pageSize
	return fmt.Sprintf(" LIMIT %d OFFSET %d ", pageSize, offset)
}
