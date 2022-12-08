package query

import (
	"reflect"

	"github.com/maniizu3110/attendance/rpc/project/external/mycache"
	"gorm.io/gorm"
)

func ApplyJoins[T any](model *T, joins []string, db *gorm.DB, cache mycache.Cache) (*gorm.DB, error) {
	joinableTables := []string{}
	for _, join := range joins {
		t := reflect.TypeOf(model).Elem().Name()
		if cache.IsJoinable(t, join) {
			joinableTables = append(joinableTables, join)
			continue
		}
		CheckExistFieldOnStruct(model, join)
		joinableTables = append(joinableTables, join)
	}
	for _, join := range joinableTables {
		db = db.Joins(join)
	}
	return db, nil
}
