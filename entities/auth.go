package entities

import (
	"github.com/scylladb/gocqlx/v2/table"
	"time"
)

type Auth struct {
	Phone     string     `db:"phone,part"`
	Code      string     `db:"code"`
	CreatedAt *time.Time `db:"created_at"`
}

func (a Auth) GetTableMetaData() table.Metadata {
	//columns, partKeys, sortKeys := getKeys(a)
	return table.Metadata{
		Name:    "auth",
		Columns: []string{"phone", "code", "created_at"},
		PartKey: []string{"phone"},
		SortKey: nil,
	}
}

func (a Auth) UpdatableKeys() []string {
	return nil
}

/*
CREATE TABLE IF NOT EXISTS auth
(
    phone      varchar PRIMARY KEY,
    Code       varchar,
    created_at timestamp
);
*/
