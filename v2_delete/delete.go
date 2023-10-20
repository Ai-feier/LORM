//go:build del
package v2_delete

import (
	"reflect"
	"strings"
)

type Deleter[T any] struct {
	Builder
	tableName string
	where []Predicate
}

func (d *Deleter[T]) Build() (*Query, error) {
	if d.sb == nil {
		d.sb = &strings.Builder{}
	}
	_, _ = d.sb.WriteString("DELETE FROM ")
	if d.tableName == "" {
		var t T 
		d.sb.WriteByte('`')
		d.sb.WriteString(reflect.TypeOf(t).Name())
		d.sb.WriteByte('`')
	} else {
		d.sb.WriteString(d.tableName)
	}
	
	if len(d.where) > 0 {
		// 构建 where 
		d.sb.WriteString(" WHERE ")
		p := d.where[0] 
		for i:=1;i< len(d.where);i++ {
			p.And(d.where[i])
		}
		if err := d.buildExpression(p); err != nil {
			return nil, err
		}
	}
	d.sb.WriteByte(';')
	return &Query{
		SQL: d.sb.String(),
		Args: d.args,
	}, nil
}

// From accepts model definition
func (d *Deleter[T]) From(table string) *Deleter[T] {
	d.tableName = table
	return d
}

// Where accepts predicates
func (d *Deleter[T]) Where(predicates ...Predicate) *Deleter[T] {
	d.where = predicates
	return d
}

//func (d *Deleter[T]) buildExpress(e Expression) error {
//	if e == nil {
//		return nil
//	}
//	switch exp := e.(type) {
//	case Column:
//		d.sb.WriteByte('`')
//		d.sb.WriteString(exp.name)
//		d.sb.WriteByte('`')
//	case value:
//		d.sb.WriteByte('?')
//		d.args = append(d.args, exp.val)
//	case Predicate:
//		_, lp := exp.left.(Predicate)
//		if lp {
//			d.sb.WriteByte('(')
//		}
//		if err := d.buildExpress(exp.left); err != nil {
//			return err
//		}
//		if lp {
//			d.sb.WriteByte(')')
//		}
//		
//		d.sb.WriteByte(' ')
//		d.sb.WriteString(exp.op.String())
//		d.sb.WriteByte(' ')
//		
//		_, rp := exp.right.(Predicate)
//		if rp {
//			d.sb.WriteByte('(')
//		}
//		if err := d.buildExpress(exp.right); err != nil {
//			return err
//		}
//		if rp {
//			d.sb.WriteByte(')')
//		}
//	default:
//		return fmt.Errorf("orm: 不支持的表达式 %v", exp)
//	}
//	return nil
//}
