package models

import (
	"strings"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

type ObjectKey string

func (k ObjectKey) String() string {
	return string(k)
}

func CreateObjectKey(filename, entityName string) ObjectKey {
	id := uuid.New()

	elements := lo.Compact([]string{entityName, filename, id.String()})

	return ObjectKey(strings.Join(elements, "/"))
}
