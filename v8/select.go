//go:build v8
package orm

import (
	"LORM/v8/internal/errs"
	"context"
	"strings"
)

// Selector 用于构造 SELECT 语句
type Selector[T any] struct {
	Builder
	table string
	db *DB
	where []Predicate
	column []Selectable
}

// From 指定表名，如果是空字符串，那么将会使用默认表名
func (s *Selector[T]) From(tbl string) *Selector[T] {
	s.table = tbl
	return s
}

func (s *Selector[T]) Build() (*Query, error) {
	var t T
	var err error
	s.Builder.model, err = s.db.r.Get(&t)
	if s.sb == nil {
		s.sb = &strings.Builder{}
	}
	s.sb.WriteString("SELECT ")
	if err = s.buildColumn(); err != nil {
		return nil, err
	}
	s.sb.WriteString(" FROM ")
	if s.table == "" {
		s.sb.WriteByte('`')
		s.sb.WriteString(s.model.TableName)
		s.sb.WriteByte('`')
	} else {
		s.sb.WriteString(s.table)
	}

	// 构造 WHERE
	if len(s.where) > 0 {
		// 类似这种可有可无的部分，都要在前面加一个空格
		s.sb.WriteString(" WHERE ")
		p := s.where[0]
		for i := 1; i < len(s.where); i++ {
			p = p.And(s.where[i])
		}
		if err := s.buildExpression(p); err != nil {
			return nil, err
		}
	}
	s.sb.WriteString(";")
	return &Query{
		SQL:  s.sb.String(),
		Args: s.args,
	}, nil
}

func (s *Selector[T]) buildColumn() error {
	if len(s.column) == 0 {
		s.sb.WriteByte('*')
		return nil
	}
	for i, c := range s.column {
		if i > 0 {
			s.sb.WriteByte(',')
		}
		switch val := c.(type) {
		case Column:
			s.sb.WriteByte('`')
			fd, ok := s.model.FieldMap[val.name]
			if !ok {
				return errs.NewErrUnknownColumn(val.name)
			}
			s.sb.WriteString(fd.ColName)
			s.sb.WriteByte('`')
			s.buildAs(val.alias)
		case Aggregate:
			s.sb.WriteString(val.fn)
			s.sb.WriteByte('(')
			fd, ok := s.model.FieldMap[val.arg]
			if !ok {
				return errs.NewErrUnknownField(val.arg)
			}
			s.sb.WriteString(fd.ColName)
			s.sb.WriteByte(')')
			s.buildAs(val.alias)
		case RawExpr:
			s.sb.WriteString(val.raw)
			if len(val.args) != 0 { 
				s.args = append(s.args, val.args...)
			}
		default:
			return errs.NewErrUnsupportedSelectable(c)
		}
	}
	return nil
}

// Where 用于构造 WHERE 查询条件。如果 ps 长度为 0，那么不会构造 WHERE 部分
func (s *Selector[T]) Where(ps ...Predicate) *Selector[T] {
	s.where = ps
	return s
}

func (s *Selector[T]) buildAs(alias string) {
	if alias != "" {
		s.sb.WriteString(" AS ")
		s.sb.WriteByte('`')
		s.sb.WriteString(alias)
		s.sb.WriteByte('`')
	}
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	q, err := s.Build()
	if err != nil {
		return nil, err
	}

	// s.db 是我们定义的 DB
	// s.db.db 则是 sql.DB
	// 使用 QueryContext，从而和 GetMulti 能够复用处理结果集的代码
	rows, err := s.db.db.QueryContext(ctx, q.SQL, q.Args...)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, ErrNoRows
	}
	
	tp := new(T)
	meta, err := s.db.r.Get(tp)
	val := s.db.valCreator(tp, meta)
	err = val.SetColumns(rows)
	return tp, err
}

func NewSelector[T any](db *DB) *Selector[T] {
	return &Selector[T]{
		db: db,
	}
}

type Selectable interface {
	selectable()
}