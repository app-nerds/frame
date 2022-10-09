package frame

import (
	"net/http"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func (fa *FrameApplication) Database(dst ...interface{}) *FrameApplication {
	var (
		err error
	)

	if fa.DB, err = gorm.Open(postgres.Open(fa.Config.DSN), &gorm.Config{}); err != nil {
		fa.Logger.WithError(err).Fatal("unable to connect to the database")
	}

	dst = append(dst, &MembersStatus{}, &Member{})
	_ = fa.DB.AutoMigrate(dst...)

	if err = fa.seedDataMemberStatuses(); err != nil {
		fa.Logger.WithError(err).Fatal("error seeding database...")
	}

	return fa
}

func (fa *FrameApplication) paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var (
			err        error
			pageString string
			page       int
		)

		if r.Method == http.MethodGet {
			pageString = r.URL.Query().Get("page")
		}

		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			pageString = r.FormValue("page")
		}

		if page, err = strconv.Atoi(pageString); err != nil {
			page = 1
		}

		if page == 0 {
			page = 1
		}

		offset := (page - 1) * fa.pageSize
		return db.Offset(offset).Limit(fa.pageSize)
	}
}
