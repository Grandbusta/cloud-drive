package models

import (
	"time"

	"github.com/Grandbusta/cloud-drive/config"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Resource struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"not null;" json:"name"`
	ResourceType string    `gorm:"not null;" json:"resource_type"`
	ParentID     string    `json:"parent_id"`
	StorageUrl   string    `json:"storage_url"`
	FileExt      string    `json:"file_ext"`
	Icon         string    `json:"icon"`
	AccessType   string    `gorm:"not null" json:"access_type"`
	UserID       string    `json:"user_id"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (r *Resource) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.NewV4().String()
	return
}

// func (r *Resource) CreateRootResource(db *gorm.DB) (*Resource, error) {
// 	r.FileExt = ".root"
// 	r.Name = "My drive"
// 	if err := db.Debug().Create(&r).Error; err != nil {
// 		return &Resource{}, err
// 	}
// 	return r, nil
// }

func (r *Resource) CreateResource(db *gorm.DB) (*Resource, error) {
	r.Prepare()
	if err := db.Debug().Create(&r).Error; err != nil {
		return &Resource{}, err
	}
	return r, nil
}

func (r *Resource) Prepare() {
	r.AccessType = config.ACCESS_TYPE_PRIVATE
}
