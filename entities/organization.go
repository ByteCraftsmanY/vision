package entities

import (
	"github.com/scylladb/gocqlx/v2/table"
	"vision/types"
)

type Organization struct {
	Base
	Name             string    `json:"name,omitempty" db:"name"`
	Contact          string    `json:"contact,omitempty" db:"contact"`
	Type             string    `json:"type,omitempty" db:"type"`
	BillAmount       int64     `json:"bill_amount,omitempty" db:"bill_amount"`
	Address          string    `json:"address,omitempty" db:"address"`
	CCTVCount        int64     `json:"cctv_count,omitempty" db:"cctv_count"`
	AssociatedUserID *types.ID `json:"associated_user_id,omitempty" db:"associated_user_id"`
	UserCount        int64     `json:"user_count,omitempty" db:"user_count"`
}

func (o Organization) GetTableMetaData() table.Metadata {
	columns, partKeys, sortKeys := getKeys(o)
	return table.Metadata{
		Name:    "organization",
		Columns: columns,
		PartKey: partKeys,
		SortKey: sortKeys,
	}
}

func (o Organization) UpdatableKeys() []string {
	return []string{"name", "address", "type", "associated_user_id"}
}
