package query

import (
	"errors"
	"reflect"
	"strings"
	"time"

	"gorm.io/gorm"
)

// membersからOrderクエリの配列を生成
func CreateOrderArray[T any](model T, members []string) ([]string, error) {
	typeOfModel := GetElementType(reflect.TypeOf(model))
	var orders []string
	for i := range members {
		order, err := CreateOrder(typeOfModel, members[i])
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

// 単一のmemberを取得して解析し，適切にOrderクエリを生成
func CreateOrder(typeOfModel reflect.Type, member string) (string, error) {
	var key string
	var dir string
	if strings.HasPrefix(member, "-") { // memberの接頭辞が"-"なら降順（desc）
		key = strings.TrimLeft(member, "-")
		key = strings.TrimSpace(key)
		dir = "desc"
	} else { // memberの接頭辞が"-"でないなら昇順（asc）
		key = strings.TrimLeft(member, " ")
		key = strings.TrimSpace(key)
		dir = "asc"
	}
	field, found := FindFieldByNameDeep(typeOfModel, key)
	if found && isSupportedType(field.Type) {
		return ToDBName(key) + " " + dir, nil
	}
	return "", errors.New(key + "がデータに含まれない為，整列できません．")
}

// Orderクエリをデータベースに適用
func ApplyOrders[T any](model T, orderFields []string, db *gorm.DB) (*gorm.DB, error) {
	orders, err := CreateOrderArray(model, orderFields)
	if err != nil {
		return nil, err
	}

	// Orderを以下のように指定
	// 		db.Order(orders[0]).Order(orders[1]).Order(orders[2])....
	// https://gorm.io/ja_JP/docs/query.html#Order
	for i := range orders {
		db = db.Order(orders[i])
	}
	return db, nil
}

func isSupportedType(t reflect.Type) bool {
	if t.Kind() == reflect.Ptr {
		return isSupportedType(t.Elem())
	}
	supportedKinds := []reflect.Kind{
		// from type.go
		reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32,
		reflect.Float64,
		reflect.String,
	}
	supportedTypes := []reflect.Type{
		reflect.TypeOf(time.Time{}),
		reflect.TypeOf(gorm.DeletedAt{}),
	}
	for i := range supportedKinds {
		if t.Kind() == supportedKinds[i] {
			return true
		}
	}
	for i := range supportedTypes {
		if t == supportedTypes[i] {
			return true
		}
	}
	return false
}
