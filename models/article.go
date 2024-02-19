package models

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ArticleStatus string

const (
	PUBLISHED ArticleStatus = "PUBLISHED"
	DRAFTED   ArticleStatus = "DRAFTED"
)

// Article is model to implement fields in database.
type Article struct {
	ID        int              `gorm:"primaryKey" json:"id"`
	UserID    int              `json:"user_id"`
	Title     string           `gorm:"size:50" json:"title"`
	Slug      string           `gorm:"size:60;uniqueIndex" json:"slug"`
	Content   string           `gorm:"type:text" json:"content"`
	Views     int              `gorm:"default:0" json:"views"`
	Status    ArticleStatus    `gorm:"default:DRAFTED" json:"status"`
	UpdatedAt time.Time        `json:"updated_at"`
	CreatedAt time.Time        `json:"created_at"`
	DeletedAt gorm.DeletedAt   `gorm:"index" json:"deleted_at"`
	Comments  []ArticleComment `gorm:"foreignKey:ArticleID" json:"comments"`
}

// ArticleModel struct to article model.
type ArticleModel struct {
	db *gorm.DB
}

// NewArticleModel is function to run article model.
func NewArticleModel(db *gorm.DB) *ArticleModel {
	return &ArticleModel{
		db: db,
	}
}

// CreateNewArticle is function to create new article.
func (a *ArticleModel) CreateNewArticle(article *Article) error {
	return a.db.Create(article).Error
}

// GetAllArticle is function to get all article.
func (a *ArticleModel) GetAllArticle(status ArticleStatus, limit, offset int) []Article {
	var articles []Article
	a.db.Model(&Article{}).Where("status = ?", status).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "updated_at"}, Desc: true}).
		Preload("Comments").Limit(limit).Offset(offset).Find(&articles)

	return articles
}

// GetArticleBySlug is function to get article by given slug.
func (a *ArticleModel) GetArticleBySlug(slug string) (Article, error) {
	var article Article
	err := a.db.Model(&Article{}).Where("slug = ?", slug).Preload("Comments").First(&article).Error
	return article, err
}
