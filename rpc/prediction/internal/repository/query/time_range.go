package query

import (
	"reflect"
	"time"

	"gorm.io/gorm"
)

// OptionWhere v がゼロ値でなければ where する
// warn : colName は インジェクション対策していない
func OptionWhere(q *gorm.DB, v any, colName string) *gorm.DB {
	if v != reflect.Zero(reflect.TypeOf(v)).Interface() {
		q = q.Where(colName+" = ?", v)
	}
	return q
}

// RangeLooseWhere [start, end) に少しでもかぶる予定を where する
// start is zero or end is zero なら何もしない
// warn : startColName, endColName は インジェクション対策していない
func RangeLooseWhere(q *gorm.DB, start time.Time, end time.Time, startColName string, endColName string) *gorm.DB {
	if !start.IsZero() && !end.IsZero() {
		q = q.Where(startColName+" < ? AND ? < "+endColName, end, start)
	}
	return q
}

// InRangeWhere [start, end) に入るデータを where
// start is zero なら start だけ無視的に動く
// warn : colName は インジェクション対策していない
func InRangeWhere(q *gorm.DB, start time.Time, end time.Time, colName string) *gorm.DB {
	if !start.IsZero() {
		q = q.Where("? <= "+colName, start)
	}
	if !end.IsZero() {
		q = q.Where(colName+" < ?", end)
	}
	return q
}
