package entities

import (
	"github.com/scylladb/gocqlx/v2/table"
	"vision/types"
)

type User struct {
	Base
	Name            string      `json:"name,omitempty" db:"name"`
	Phone           string      `json:"phone,omitempty" db:"phone"`
	Email           string      `json:"email,omitempty" db:"email"`
	Password        string      `json:"-" db:"password"`
	OrganizationIDs []*types.ID `json:"organization_ids,omitempty" db:"organization_ids"`
}

func (u User) GetTableMetaData() table.Metadata {
	columns, partKeys, sortKeys := getKeys(u)
	return table.Metadata{
		Name:    "user",
		Columns: columns,
		PartKey: partKeys,
		SortKey: sortKeys,
	}
}

func (u User) UpdatableKeys() []string {
	return nil
}
