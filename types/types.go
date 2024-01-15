package types

import (
	"github.com/gocql/gocql"
)

type ID = gocql.UUID

func CreateNewID() *ID {
	id := gocql.TimeUUID()
	return &id
}

func ParseID(str string) (*ID, error) {
	id, err := gocql.ParseUUID(str)
	return &id, err
}
