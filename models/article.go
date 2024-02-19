package models

import (
	"time"

	"gorm.io/gorm"
)

type ArticleStatus string

const (
	PUBLISHED ArticleStatus = "PUBLISHED"
	DRAFTED   ArticleStatus = "DRAFTED"
)

type Article struct {
	ID        int            `gorm:"primaryKey" json:"id"`
	UserID    int            `json:"user_id"`
	Title     string         `gorm:"size:50" json:"title"`
	Slug      string         `gorm:"size:60;uniqueIndex" json:"slug"`
	Content   string         `gorm:"type:text" json:"content"`
	Views     int            `gorm:"default:0" json:"views"`
	Status    ArticleStatus  `gorm:"default:DRAFTED" json:"status"`
	UpdatedAt time.Time      `json:"updated_at"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
