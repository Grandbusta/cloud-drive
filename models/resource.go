package models

import "gorm.io/gorm"

type Resource struct {
	Name         string `gorm:"not null;" json:"name"`
	ResourceType string `gorm:"not null;" json:"resource_type"`
	ParentID     string `json:"parent_id"`
	StorageUrl   string `json:"storage_url"`
	FileExt      string `json:"file_ext"`
	Icon         string `json:"icon"`
	AccessType   string `gorm:"not null" json:"access_type"`
	UserID       uint   `json:"user_id"`
	gorm.Model
}
