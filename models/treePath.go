package models

import "gorm.io/gorm"

type TreePath struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Ancestor   string `gorm:"ancestor" json:"ancestor"`
	Descendant string `gorm:"descendant" json:"descendant"`
	Depth      uint   `gormm:"depth" json:"depth"`
}

// Insert a root node
func (t *TreePath) InsertRoot(db *gorm.DB) error {
	return db.Exec("INSERT INTO tree_paths(ancestor,descendant,depth) VALUES (?,?,0)", t.Ancestor, t.Descendant).Error
}

// add a new child to a node
func (t *TreePath) InsertDescendant(db *gorm.DB) error {
	err := db.Exec("INSERT INTO tree_paths(ancestor,descendant,depth) VALUES (?,?,0)", t.Descendant, t.Descendant).Error
	if err != nil {
		return err
	}
	return t.LinkDescendant(db)
}

// link a child node to a parent
func (t *TreePath) LinkDescendant(db *gorm.DB) error {
	return db.Exec(`INSERT INTO tree_paths(ancestor,descendant,depth) SELECT a.ancestor, d.descendant, a.depth + d.depth + 1 FROM tree_paths a, tree_paths d WHERE a.descendant = ? AND d.ancestor = ?`, t.Ancestor, t.Descendant).Error
}

// delete sub-tree
func (t *TreePath) DeleteDescendants(db *gorm.DB) error {
	return db.Exec("DELETE FROM tree WHERE parent = ? AND child <> ?", t.Ancestor, t.Ancestor).Error
}
