package utils

import (
	"github.com/gofrs/uuid/v5"
)

func ConvertUUIDListToStringList(uuids []uuid.UUID) []string {
	strList := make([]string, len(uuids))

	for i, u := range uuids {
		strList[i] = u.String()
	}
	return strList
}
