package http

import "github.com/yarikTri/archipelago-notes-api/internal/models"

type methodName uint8

const (
	getMethodName methodName = iota
	updateMethodName
	deleteMethodName
	setAccessMethodName
)

func (mn *methodName) String() string {
	switch *mn {
	case getMethodName:
		return "get"
	case updateMethodName:
		return "update"
	case deleteMethodName:
		return "delete"
	case setAccessMethodName:
		return "set_access"
	}

	return ""
}

var methodsAccessMap = map[methodName][]models.NoteAccess{
	getMethodName:       {models.ReadNoteAccess, models.WriteNoteAccess, models.ModifyNoteAccess, models.ManageAccessNoteAccess},
	updateMethodName:    {models.ModifyNoteAccess, models.ManageAccessNoteAccess},
	deleteMethodName:    {models.ModifyNoteAccess, models.ManageAccessNoteAccess},
	setAccessMethodName: {models.ManageAccessNoteAccess},
}

func getAllowedMethods(access models.NoteAccess) []string {
	allowedMethods := make([]string, 0)
	
	for method, accesses := range methodsAccessMap {
		for _, a := range accesses {
			if a == access {
				allowedMethods = append(allowedMethods, method.String())
				break
			}
		}
	}

	return allowedMethods
}
