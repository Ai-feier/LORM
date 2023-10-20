//go:build v3
package orm

import (
	"LORM/v3/internal/errs"
	"reflect"
	"unicode"
)

type model struct {
	// tableName 结构体对应的表名
	tableName string
	fieldMap map[string]*field 
}

// field 字段
type field struct {
	colName string
}

// 接受一级结构体指针
func parseModel(val any) (*model, error) {
	typ := reflect.TypeOf(val)
	if typ.Kind() != reflect.Pointer || 
		typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}
	typ = typ.Elem()
	
	
	// 获取字段数量
	numField := typ.NumField()
	fds := make(map[string]*field, numField)
	for i := 0; i < numField; i++ {
		fdTyp := typ.Field(i)
		fds[fdTyp.Name] = &field{
			colName: underscoreName(fdTyp.Name),
		}
	}
	return &model{
		tableName: underscoreName(typ.Name()),
		fieldMap:  fds,
	}, nil
	
}

// underscoreName 驼峰转字符串命名
func underscoreName(tableName string) string {
	var buf []byte
	for i, v := range tableName {
		if unicode.IsUpper(v) {
			if i != 0 {
				buf = append(buf, '_')
			}
			buf = append(buf, byte(unicode.ToLower(v)))
		} else {
			buf = append(buf, byte(v))
		}
	}
	return string(buf)
} 