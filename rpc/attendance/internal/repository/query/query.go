package query

import (
	"errors"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

const (
	NestedQuerySeparator = "."
)

// 条件(condition)を適切なqueryに変換．
// 変換できない場合はok=false．
func ConvertConditionToQuery(t reflect.Type, condition string, substr string) (key string, vals []string, ok bool) {
	key, values, ok := strings.Cut(condition, substr) // 条件(condifion)が指定文字列(substr)で区切れるか判定
	if !ok {
		return
	}

	// joinしたものを検索する場合のための処理　例) User.ID => ID
	t2, key, isNestedQuery := strings.Cut(key, NestedQuerySeparator)
	if !isNestedQuery {
		key = t2
	}

	for _, str := range []string{"<=", ">=", "!=", "=", "<", ">"} { // ["< =", "<< =", "> =", ">> ="]みたいなパターンがないか
		_, _, recutable := strings.Cut(key, str)
		if recutable {
			ok = false
			return
		}
		_, _, recutable = strings.Cut(values, str)
		if recutable {
			ok = false
			return
		}
	}

	key = strings.TrimSpace(key) // 空文字列をcutしてkeyのみ抽出 ex.) "Key   " -> "Key"
	if isNestedQuery {
		field, found := FindFieldByNameDeep(t, t2)
		if !found {
			ok = false
			return
		}
		searchedType := toRawStruct(field.Type)
		_, found = FindFieldByNameDeep(searchedType, key)
		if !found {
			ok = false
			return
		}
		key = ToDBName(key) // keyをdb用の名前に変換
		key = t2 + NestedQuerySeparator + key
	} else {
		field, found := FindFieldByNameDeep(t, key)
		if !found || (found && !isSupportedType(field.Type)) { // fieldを持たない場合，またはfieldを持ち且つそのfieldがサポートされていないタイプなら検索対象外
			ok = false
			return
		}
		key = ToDBName(key) // keyをdb用の名前に変換
	}
	for _, value := range strings.Split(values, ",") { // valueを文字列配列型に変換
		vals = append(vals, strings.TrimSpace(value))
	}
	ok = len(vals) > 0
	return
}

// 単一の条件(condition)を取得して，適切なqueryに変換
func CreateQuery(t reflect.Type, condition string) ([]any, error) {
	if key, vals, ok := ConvertConditionToQuery(t, condition, " not in "); ok {
		return []any{key + " NOT IN (?)", vals}, nil
	} else if key, vals, ok := ConvertConditionToQuery(t, condition, " in "); ok {
		return []any{key + " IN (?)", vals}, nil
	} else if key, vals, ok := ConvertConditionToQuery(t, condition, " includes "); ok && len(vals) == 1 {
		return []any{key + " LIKE ?", "%" + vals[0] + "%"}, nil
	}
	for _, substr := range []string{"<=", ">=", "!=", "=", "<", ">"} {
		if key, vals, ok := ConvertConditionToQuery(t, condition, substr); ok && len(vals) == 1 {
			return []any{key + " " + substr + " ?", vals[0]}, nil
		}
	}
	if key, vals, ok := ConvertConditionToQuery(t, condition, " is "); ok && len(vals) == 1 {
		switch vals[0] {
		case "NULL":
			return []any{key + " is NULL"}, nil
		case "NOTNULL":
			return []any{key + " is NOT NULL"}, nil
		}
	}
	return nil, errors.New(condition + " は無効なクエリです")
}

// 条件(conditions)をqueryに変換して多次元配列形式で返却
func CreateQueryArray[T any](model T, fetchConfig FetchConfig) ([][]any, error) {
	var queries [][]any
	conditions := fetchConfig.Query
	t := GetElementType(reflect.TypeOf(model))
	for _, condition := range conditions { // 各条件(condition)からqueryを作成
		if strings.Contains(condition, ".") {
			CheckIsQuerySupported(fetchConfig.Joins, condition)
		}
		query, err := CreateQuery(t, condition)
		if err != nil {
			return nil, err
		}
		queries = append(queries, query)
	}
	return queries, nil
}

// Queryクエリをデータベースに適用
func ApplyQueries[T any](model T, fetchConfig FetchConfig, db *gorm.DB) (*gorm.DB, error) {
	queries, err := CreateQueryArray(model, fetchConfig)
	if err != nil {
		return nil, err
	}
	for i := range queries {
		db = db.Where(queries[i][0], queries[i][1:]...)
	}
	return db, nil
}

func CheckIsQuerySupported(joins []string, condition string) error {
	key, _, _ := strings.Cut(condition, ".")
	for _, join := range joins {
		if strings.Contains(join, key) {
			return nil
		}
	}
	return errors.New("invalid query field: " + condition)
}
