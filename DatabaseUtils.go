package frame

import "fmt"

func GetDBPaging(page int, pageSize int) string {
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * pageSize
	return fmt.Sprintf(" LIMIT %d OFFSET %d ", pageSize, offset)
}
