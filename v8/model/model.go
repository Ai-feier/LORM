//go:build v8

package model

import (
	"LORM/v8/internal/errs"
	"reflect"
)

type ModelOpt func(model *Model) error

type Model struct {
	// TableName 结构体对应的表名
	TableName string
	FieldMap  map[string]*Field
	ColumnMap map[string]*Field
}

func WithTableName(name string) ModelOpt {
	return func(m *Model) error {
		m.TableName = name
		return nil
	}
}

func WithColumnName(field, colName string) ModelOpt {
	return func(m *Model) error {
		fd, ok := m.FieldMap[field]
		if !ok {
			return errs.NewErrUnknownField(field)
		}
		// 注意，这里我们根本没有检测 ColName 会不会是空字符串
		// 因为正常情况下，用户都不会写错
		// 即便写错了，也很容易在测试中发现
		fd.ColName = colName
		return nil
	}
}

// Field 字段
type Field struct {
	ColName string
	GoName string
	Type reflect.Type
	// Offset 相对于对象起始地址的字段偏移量
	Offset uintptr
}

// 我们支持的全部标签上的 key 都放在这里
// 方便用户查找，和我们后期维护
const (
	tagKeyColumn = "column"
)

// 用户自定义一些模型信息的接口，集中放在这里
// 方便用户查找和我们后期维护

// TableName 用户实现这个接口来返回自定义的表名
type TableName interface {
	TableName() string
}
