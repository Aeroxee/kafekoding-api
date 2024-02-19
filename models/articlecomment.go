package models

import (
	"time"

	"gorm.io/gorm"
)

// ArticleComment
type ArticleComment struct {
	ID        int            `gorm:"primaryKey" json:"id"`
	ArticleID int            `json:"article_id"`
	UserID    int            `json:"user_id"`
	Text      string         `gorm:"type:text" json:"text"`
	UpdatedAt time.Time      `json:"updated_at"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
