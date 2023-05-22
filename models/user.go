package models

import "gorm.io/gorm"

type User struct {
	Email          string     `gorm:"unique;not null" json:"email"`
	Password       string     `gorm:"not null;" json:"password"`
	RootResourceID uint       `json:"root_resource_id"`
	Resources      []Resource `json:"resources"`
	gorm.Model
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	if err := db.Debug().Create(&u).Error; err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) FindUserByEmail(db *gorm.DB) (*User, error) {
	if err := db.Debug().Where("email=?", u.Email).First(&u).Error; err != nil {
		return &User{}, nil
	}
	return u, nil
}
