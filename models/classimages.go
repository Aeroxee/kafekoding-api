package models

import (
	"time"

	"gorm.io/gorm"
)

type ClassImage struct {
	ID        int            `gorm:"primaryKey" json:"id"`
	ClassID   int            `json:"class_id"`
	Image     string         `gorm:"size:255" json:"image"`
	Caption   *string        `gorm:"size:255" json:"caption"`
	UpdatedAt time.Time      `json:"updated_at"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
