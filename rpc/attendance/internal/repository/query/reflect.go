package query

import (
	"reflect"
	"strings"
)

// 指定の型(t)に指定の名前(name)のフィールドが存在するかの判定と，含まれた場合はその情報を返却（埋め込まれた構造体も対象）
func FindFieldByNameDeep(t reflect.Type, name string) (reflect.StructField, bool) {
	f, found := t.FieldByName(name)
	if found {
		return f, true
	}
	for i := 0; i < t.NumField(); i++ {
		switch t.Field(i).Type.Kind() {
		case reflect.Struct:
			if t.Field(i).Anonymous {
				f, found = FindFieldByNameDeep(t.Field(i).Type, name)
				if found {
					return f, found
				}
			}
		}
	}
	return reflect.StructField{}, false
}

// T型のデータの型情報を取得
func GetElementTypeDeep[T any](model T) reflect.Type {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	return t
}

// T型のデータに入っている値情報の取得
func GetElementValue[T any](model T) reflect.Value {
	vs := reflect.ValueOf(model)
	if vs.Kind() == reflect.Ptr {
		vs = vs.Elem()
	}
	return vs
}

// 型情報を取得
func GetElementType(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// field が tag: value を持つかチェック
func FindValueFromTag[T any](field reflect.StructField, tag string, value string) bool {
	tags := strings.Split(field.Tag.Get(tag), ",")
	for i := range tags {
		if strings.TrimSpace(tags[i]) == value {
			return true
		}
	}
	return false
}

func NewInstance(i any) any {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return reflect.New(t).Interface()
}

func NewSliceOf(i any) any {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	slice := reflect.MakeSlice(reflect.SliceOf(t), 0, 0)
	ptr := reflect.New(slice.Type())
	ptr.Elem().Set(slice)
	return ptr.Interface()
}

// copy
func ShallowCopy(m any) any {
	t := reflect.TypeOf(m)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	clone := reflect.New(t)
	clone.Elem().Set(reflect.Indirect(reflect.ValueOf(m)))
	return clone.Interface()
}
