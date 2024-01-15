package daos

import (
	"github.com/scylladb/gocqlx/v2/table"
	"vision/types"
)

type Product struct {
	Base
	Username       string    `json:"username,omitempty"`
	Password       string    `json:"-"`
	URL            string    `json:"url,omitempty"`
	OrganizationID *types.ID `json:"organization_id,omitempty"`
}

func (p Product) GetTableMetaData() table.Metadata {
	columns, partKeys, sortKeys := getKeys(p)
	return table.Metadata{
		Name:    "organization",
		Columns: columns,
		PartKey: partKeys,
		SortKey: sortKeys,
	}
}

func (p Product) UpdatableKeys() []string {
	return nil
}
