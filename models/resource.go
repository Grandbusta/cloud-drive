package models

import "gorm.io/gorm"

type Resource struct {
	gorm.Model
	Name         string `json:"name"`
	ResourceType string `json:"resource_type"`
	ParentID     string `json:"parent_id"`
	StorageUrl   string `json:"storage_url"`
	FileExt      string `json:"file_ext"`
	Icon         string `json:"icon"`
	AccessType   string `json:"access_type"`
	UserID       uint   `json:"user_id"`
}
