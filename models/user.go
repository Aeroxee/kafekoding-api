package models

import (
	"time"

	"github.com/Aeroxee/kafekoding-api/auth"
)

type UserType int8

const (
	ADMIN UserType = iota
	MEMBER
)

type User struct {
	ID           int       `gorm:"primaryKey" json:"id,omitempty"`
	FirstName    string    `gorm:"size:50" json:"first_name,omitempty"`
	LastName     string    `gorm:"size:50" json:"last_name,omitempty"`
	Username     string    `gorm:"size:50;uniqueIndex" json:"username,omitempty"`
	Email        string    `gorm:"size:50;uniqueIndex" json:"email,omitempty"`
	Avatar       *string   `gorm:"size:255" json:"avatar,omitempty"`
	Password     string    `gorm:"size:128" json:"-"`
	IsActive     bool      `gorm:"default:false" json:"is_active,omitempty"`
	IsLogined    bool      `gorm:"default:false" json:"is_logined,omitempty"`
	Type         UserType  `gorm:"default:1" json:"type"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
	DateJoined   time.Time `gorm:"autoCreateTime" json:"date_joined,omitempty"`
	Articles     []Article `gorm:"foreignKey:UserID" json:"articles,omitempty"`
	ClassMentors []*Class  `gorm:"many2many:classes_user_mentor" json:"class_mentors,omitempty"`
	ClassMembers []*Class  `gorm:"many2many:classes_user_member" json:"class_members,omitempty"`
}

func CreateNewUser(user *User) error {
	user.Password = auth.EncryptionPassword(user.Password)
	return DB().Create(user).Error
}

func GetUserByID(id int) (User, error) {
	var user User
	err := DB().Model(&User{}).Where("id = ?", id).Preload("Articles").Preload("ClassMentors").
		Preload("ClassMembers").First(&user).Error
	return user, err
}

func GetUserByUsername(username string) (User, error) {
	var user User
	err := DB().Model(&User{}).Where("username = ?", username).Preload("Articles").Preload("ClassMentors").
		Preload("ClassMembers").First(&user).Error
	return user, err
}

func GetUserByEmail(email string) (User, error) {
	var user User
	err := DB().Model(&User{}).Where("email = ?", email).Preload("Articles").Preload("ClassMentors").
		Preload("ClassMembers").First(&user).Error
	return user, err
}
