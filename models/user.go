package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        string     `gorm:"primaryKey" json:"id"`
	Email     string     `gorm:"unique;not null" json:"email"`
	Password  string     `gorm:"not null;" json:"password"`
	Resources []Resource `json:"resources"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type CreateUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PublicUser struct {
	ID        string     `json:"id"`
	Email     string     `json:"email"`
	Resources []Resource `json:"resources"`
	CreatedAt time.Time  `json:"created_at"`
}

func (u *User) PublicUser() *PublicUser {
	return &PublicUser{
		ID:        u.ID,
		Email:     u.Email,
		Resources: u.Resources,
		CreatedAt: u.CreatedAt,
	}
}

func (r *User) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.NewV4().String()
	return
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	if err := db.Debug().Create(&u).Error; err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) FindUserByEmail(db *gorm.DB) (*User, error) {
	if err := db.Debug().Where("email=?", u.Email).First(&u).Error; err != nil {
		return &User{}, err
	}
	return u, nil
}
