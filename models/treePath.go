package models

import (
	"gorm.io/gorm"
)

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
func (t *TreePath) DeleteSelfWithDescendants(db *gorm.DB) error {
	return db.Exec("DELETE FROM tree_paths WHERE descendant IN (SELECT tree_paths.descendant from tree_paths where tree_paths.ancestor = ?)", t.Ancestor).Error
}

func (t *TreePath) SelectSelfWithDescendants(db *gorm.DB) ([]string, error) {
	var paths []string
	err := db.Raw("SELECT tree_paths.descendant from tree_paths where tree_paths.ancestor = ? ORDER BY depth ASC", t.Ancestor).Scan(&paths).Error
	if err != nil {
		return paths, err
	}
	return paths, nil
}

// Select direct children
func (t *TreePath) SelectChildren(db *gorm.DB) ([]PublicResource, error) {
	var resources []PublicResource
	err := db.Raw("SELECT resources.* FROM resources JOIN tree_paths t ON (resources.id = t.descendant) where t.ancestor = ? AND depth = 1", t.Ancestor).Scan(&resources).Error
	if err != nil {
		return []PublicResource{}, err
	}
	return resources, nil
}
