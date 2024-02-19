package models

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Class is struct for implement class in kafekoding.
type Class struct {
	ID          int             `gorm:"primaryKey" json:"id"`
	Title       string          `gorm:"size:50;unique" json:"title"`
	Slug        string          `gorm:"size:60;uniqueIndex" json:"slug"`
	Description string          `gorm:"type:text" json:"description"`
	Logo        *string         `gorm:"size:255" json:"logo"`
	IsActive    bool            `gorm:"default:false" json:"is_active"`
	UpdatedAt   time.Time       `json:"updated_at"`
	CreatedAt   time.Time       `json:"created_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"deleted_at"`
	Mentors     []*User         `gorm:"many2many:classes_user_mentor" json:"mentors"`
	Members     []*User         `gorm:"many2many:classes_user_member" json:"members"`
	Images      []*ClassImage   `gorm:"foreignKey:ClassID" json:"images"`
	Meetings    []*ClassMeeting `gorm:"foreignKey:ClassID" json:"meetings"`
}

func CreateNewClass(class *Class) error {
	return DB().Create(class).Error
}

func GetAllClass(is_active bool) []Class {
	var classes []Class
	DB().Model(&Class{}).Where("is_active = ?", is_active).Order(clause.OrderByColumn{
		Column: clause.Column{Name: "title"},
		Desc:   false,
	}).Preload("Mentors").Preload("Members").Preload("Images").Preload("Meetings").
		Find(&classes)

	return classes
}

func GetClassBySlug(slug string) (Class, error) {
	var class Class
	err := DB().Model(&Class{}).Where("slug = ?", slug).Preload("Mentors").Preload("Members").
		Preload("Images").Preload("Meetings").First(&class).Error
	return class, err
}
