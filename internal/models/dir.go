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
	ID       int             `json:"id"`
	Name     string          `json:"name"`
	Children []*DirTree      `json:"children"`
	Notes    []*NoteTransfer `json:"notes"`
}

func (d *Dir) ToTree(notes []*Note) *DirTree {
	var notesTransfer = make([]*NoteTransfer, 0)
	for _, note := range notes {
		notesTransfer = append(notesTransfer, note.ToTransfer())
	}

	return &DirTree{
		ID:       d.ID,
		Name:     d.Name,
		Children: make([]*DirTree, 0),
		Notes:    notesTransfer,
	}
}

func ToTree(rootID int, dirs []*Dir, notes []*Note) *DirTree {
	notesByDirID := make(map[int][]*Note)
	for _, note := range notes {
		if notesByDirID[note.DirID] == nil {
			notesByDirID[note.DirID] = make([]*Note, 0)
		}

		notesByDirID[note.DirID] = append(notesByDirID[note.DirID], note)
	}

	dirTreesByPathMap := make(map[string]*DirTree)
	dirTreesByID := make(map[int]*DirTree)

	var root *DirTree
	for _, dir := range dirs {
		keyPrefix := ""
		if dir.Path != "" {
			keyPrefix = dir.Path + "."
		}

		dirTree := dir.ToTree(notesByDirID[dir.ID])

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
