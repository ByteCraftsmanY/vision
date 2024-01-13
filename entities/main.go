package entities

import (
	"github.com/scylladb/gocqlx/v2/table"
	"reflect"
	"strings"
)

type DBModel interface {
	GetTableMetaData() table.Metadata
	UpdatableKeys() []string
}

func extractTags(model interface{}) []string {
	tags := make([]string, 0)
	structType := reflect.TypeOf(model)
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		dbTag := field.Tag.Get("db")
		if strings.TrimSpace(dbTag) == "" {
			continue
		}
		tags = append(tags, dbTag)
	}
	return tags
}

func getKeys(model DBModel) (columns, partKeys, sortKeys []string) {
	baseTags := extractTags(Base{})
	modelTags := extractTags(model)
	tags := append(baseTags, modelTags...)
	for _, tag := range tags {
		// keys: "db:column_name,part" || "db:column_name,sort"
		keys := strings.Split(tag, ",")
		column := keys[0]
		if len(keys) >= 2 {
			_type := keys[1]
			switch _type {
			case "part":
				partKeys = append(partKeys, column)
			case "sort":
				sortKeys = append(sortKeys, column)
			}
		}
		columns = append(columns, column)
	}
	return
}
