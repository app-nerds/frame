package frame

import (
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

	dst = append(dst, &Member{})
	_ = fa.DB.AutoMigrate(dst...)

	return fa
}
