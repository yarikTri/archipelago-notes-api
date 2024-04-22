package models

import (
	"strconv"
)

type Dir struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Path string `json:"subpath" db:"subpath"`
}

type DirTree struct {
	ID       int        `json:"id"`
	Name     string     `json:"name"`
	Children []*DirTree `json:"children"`
}

func (d *Dir) ToTree() *DirTree {
	return &DirTree{
		ID:       d.ID,
		Name:     d.Name,
		Children: make([]*DirTree, 0),
	}
}

func ToTree(rootID int, dirs []*Dir) *DirTree {
	dirTreesByPathMap := make(map[string]*DirTree)
	dirTreesByID := make(map[int]*DirTree)

	var root *DirTree
	for _, dir := range dirs {
		keyPrefix := ""
		if dir.Path != "" {
			keyPrefix = dir.Path + "."
		}

		dirTree := dir.ToTree()

		key := keyPrefix + strconv.Itoa(dir.ID)
		dirTreesByPathMap[key] = dirTree
		dirTreesByID[dir.ID] = dirTree

		if dir.ID == rootID {
			root = dirTree
		}
	}

	for _, dir := range dirs {
		parentDir := dirTreesByPathMap[dir.Path]
		if parentDir == nil {
			continue
		}

		parentDir.Children = append(parentDir.Children, dirTreesByID[dir.ID])
	}

	return root
}
