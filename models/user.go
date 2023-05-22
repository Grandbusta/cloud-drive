package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email          string     `gorm:"unique;not null" json:"email"`
	Password       string     `gorm:"not null;" json:"password"`
	RootResourceID Resource   `json:"root_resource_id"`
	Resources      []Resource `json:"resources"`
}
