package models

import "time"

type ClassPermissionType string

const (
	JOIN ClassPermissionType = "JOIN"
	OUT  ClassPermissionType = "OUT"
)

type ClassPermission struct {
	ID        int `gorm:"primaryKey"`
	UserID    int
	ClassID   int
	Type      ClassPermissionType `gorm:"size:4"`
	UpdatedAt time.Time
	CreatedAt time.Time
}
