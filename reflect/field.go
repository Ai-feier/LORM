package reflect

import (
	"errors"
	"reflect"
)

func IterateFields(entity any) (res map[string]any, err error) {
	if entity == nil {
		return nil, errors.New("不支持 nil")
	}
	typ := reflect.TypeOf(entity)
	val := reflect.ValueOf(entity)

	// 判斷是否為 零值指針(帶有指針類型信息的 nil)
	if val.IsZero() {
		return nil, errors.New("不支持零值")
	}
	// 對於多級指針
	for typ.Kind() == reflect.Pointer {
		// 拿到指針指向的對象
		typ = typ.Elem()
		val = val.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return nil, errors.New("仅支持 Struct")
	}

	res = map[string]any{}
	numField := val.NumField()
	for i := 0; i < numField; i++ {
		fieldType := typ.Field(i)
		// fieldValue 如果为不可导出类型会panic
		fieldValue := val.Field(i)
		if fieldType.IsExported() {
			res[fieldType.Name] = fieldValue.Interface()
		} else {
			res[fieldType.Name] = reflect.Zero(fieldType.Type).Interface()
		}

	}
	return
}

func SetField(entity any, field string, newValue any) error {
	val := reflect.ValueOf(entity)
	// 取指针的内容
	for val.Type().Kind() == reflect.Pointer {
		val = val.Elem()
	}
	fieldValue := val.FieldByName(field)
	if !fieldValue.CanSet() {
		// 不可修改
		return errors.New("不可修改")
	}
	fieldValue.Set(reflect.ValueOf(newValue))
	return nil
}
