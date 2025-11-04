package models

import "github.com/google/uuid"

type ObjectKey string

func (k ObjectKey) String() string {
	return string(k)
}

func CreateObjectKey(filename string) ObjectKey {
	id := uuid.New()

	return ObjectKey(id.String() + "/" + filename)
}
