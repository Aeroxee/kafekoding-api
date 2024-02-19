package models

import (
	"time"

	"gorm.io/gorm"
)

type ClassMeetingAttendance struct {
	ID        int            `gorm:"primaryKey" json:"id"`
	MeetingID int            `json:"meeting_id"`
	Users     []*User        `gorm:"many2many:classes_meetingattendance_user" json:"users"`
	UpdatedAt time.Time      `json:"updated_at"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
