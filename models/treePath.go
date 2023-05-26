package models

type TreePath struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Ancestor   string `gorm:"ancestor" json:"ancestor"`
	Descendant string `gorm:"descendant" json:"descendant"`
	Depth      int    `gormm:"depth" json:"depth"`
}
