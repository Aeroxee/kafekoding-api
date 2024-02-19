package models

import (
	"time"

	"gorm.io/gorm"
)

type ClassMeeting struct {
	ID          int                      `gorm:"primaryKey" json:"id"`
	ClassID     int                      `json:"class_id"`
	Title       string                   `gorm:"size:50" json:"title"`
	Slug        string                   `gorm:"size:60;uniqueIndex" json:"slug"`
	Content     string                   `gorm:"type:text" json:"content"`
	OpenedAt    time.Time                `json:"opened_at"`
	ClosedAt    time.Time                `json:"closed_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
	CreatedAt   time.Time                `json:"created_at"`
	DeletedAt   gorm.DeletedAt           `gorm:"index" json:"deleted_at"`
	Attendances []ClassMeetingAttendance `gorm:"foreignKey:MeetingID" json:"attendances"`
}
