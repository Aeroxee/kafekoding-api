package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DB() *gorm.DB {
	dsn := "root:root@tcp(127.0.0.1:3306)/kafekoding?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{}, &Class{}, &ClassMeeting{}, &ClassImage{},
		&ClassMeetingAttendance{}, &Article{}, &ArticleComment{})
	return db
}
