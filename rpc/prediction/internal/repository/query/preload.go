package query

import (
	"errors"
	"reflect"
	"strings"

	"github.com/maniizu3110/attendance/rpc/prediction/external/mycache"
	"gorm.io/gorm"
)

// Preloadを適用する
func ApplyPreloads[T any](model *T, preloadFields []string, db *gorm.DB, cache mycache.Cache) (*gorm.DB, error) {
	members := []string{}
	for _, preloadField := range preloadFields {
		typeName := reflect.TypeOf(model).Elem().Name()
		if cache.IsPreloadable(typeName, preloadField) {
			members = append(members, preloadField)
			continue
		}
		err := CheckIsPreloadableField(model, preloadField)
		if err != nil {
			return nil, err
		}
		members = append(members, preloadField)
	}
	for _, m := range members {
		db = db.Preload(m)
	}
	return db, nil
}

// フィールドが特定の構造体(ネストを含めて)に存在するかどうか
// 例: CheckIsPreloadableField(&User{}, "Profile") -> true
func CheckIsPreloadableField[T any](model T, preload string) error {
	preloadFields := splitPreloadFieldsToArray(preload)
	var (
		searchedType reflect.Type
		newModel     interface{}
		err          error
	)
	for i, preloadField := range preloadFields {
		if i == 0 {
			searchedType = GetElementTypeDeep(model)
			field, found := FindFieldByNameDeep(searchedType, preloadField)
			if found {
				searchedType = toRawStruct(field.Type)
			}
		} else {
			newModel = reflect.New(searchedType).Interface()
			err = CheckExistFieldOnStruct(newModel, preloadField)
			if err != nil {
				return err
			}
			if i == len(preloadFields)-1 {
				break
			}
			field, found := FindFieldByNameDeep(searchedType, preloadField)
			if found {
				searchedType = toRawStruct(field.Type)
			}
		}
	}
	return nil
}

func splitPreloadFieldsToArray(preloadField string) []string {
	return strings.Split(preloadField, ".")
}

// フィールドが構造体に存在するかどうか確認
func CheckExistFieldOnStruct[T any](model T, preload string) error {
	t := GetElementTypeDeep(model)
	field, found := FindFieldByNameDeep(t, preload)
	if found && !field.Anonymous { // memberと同じ名前を持つfieldが存在し，且つそれが埋め込まれたものでない場合...
		t2 := field.Type
		// ポインターの場合はポインター先の型を取得
		if t2.Kind() == reflect.Ptr {
			t2 = t2.Elem()
		}
		// unpreloadbleタグがついていないことを確認
		if t2.Kind() == reflect.Struct || t2.Kind() == reflect.Slice { // 型が構造体orスライス且つ...
			preloadable := !FindValueFromTag[T](field, "api", "unpreloadable")
			if !preloadable {
				return errors.New("invalid preload field: " + preload)
			}
			return nil
		}
	}
	return errors.New("invalid preload field: " + preload)
}

// ポインタや配列の場合は要素の型を取得
// 例: []*repository.User -> User
func toRawStruct(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr || t.Kind() == reflect.Slice {
		t = t.Elem()
		if t.Kind() == reflect.Ptr || t.Kind() == reflect.Slice {
			t = toRawStruct(t)
		}
	}
	return t
}
