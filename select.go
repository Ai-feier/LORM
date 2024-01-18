package lorm

import (
	"github.com/Ai-feier/lorm/internal/errs"
	"context"
	"database/sql"
)

// Selector 用于构造 SELECT 语句
type Selector[T any] struct {
	builder
	table   string
	db      *DB
	where   []Predicate
	having  []Predicate
	columns []Selectable
	groupBy []Column
	limit   int
	offset  int
}

// Select 选择要查询的列
func (s *Selector[T]) Select(cols ...Selectable) *Selector[T] {
	s.columns = cols
	return s
}

// From 指定表名，如果是空字符串，那么将会使用默认表名
func (s *Selector[T]) From(tbl string) *Selector[T] {
	s.table = tbl
	return s
}

// Build 构造 sql 查询语句, 底层调用 database.sql 查询数据库
func (s *Selector[T]) Build() (*Query, error) {
	var (
		t   T
		err error
	)
	s.model, err = s.db.r.Get(&t)
	if err != nil {
		return nil, err
	}
	s.sb.WriteString("SELECT ")
	if err = s.buildColumns(); err != nil {
		return nil, err
	}
	s.sb.WriteString(" FROM ")
	if s.table == "" {
		s.quote(s.model.TableName)
	} else {
		s.sb.WriteString(s.table)
	}

	// 构造 WHERE
	if len(s.where) > 0 {
		// 类似这种可有可无的部分，都要在前面加一个空格
		s.sb.WriteString(" WHERE ")
		// WHERE 是不允许用别名的
		if err = s.buildPredicates(s.where); err != nil {
			return nil, err
		}
	}

	if len(s.groupBy) > 0 {
		s.sb.WriteString(" GROUP BY ")
		for i, c := range s.groupBy {
			if i > 0 {
				s.sb.WriteByte(',')
			}
			if err = s.buildColumn(c, false); err != nil {
				return nil, err
			}
		}
	}

	if len(s.having) > 0 {
		s.sb.WriteString(" HAVING ")
		// HAVING 是可以用别名的
		if err = s.buildPredicates(s.having); err != nil {
			return nil, err
		}
	}

	if s.limit > 0 {
		s.sb.WriteString(" LIMIT ?")
		s.addArgs(s.limit)
	}

	if s.offset > 0 {
		s.sb.WriteString(" OFFSET ?")
		s.addArgs(s.offset)
	}

	s.sb.WriteString(";")
	return &Query{
		SQL:  s.sb.String(),
		Args: s.args,
	}, nil
}

func (s *Selector[T]) buildColumns() error {
	if len(s.columns) == 0 {
		s.sb.WriteByte('*')
		return nil
	}
	for i, c := range s.columns {
		if i > 0 {
			s.sb.WriteByte(',')
		}
		switch val := c.(type) {
		case Column:
			if err := s.buildColumn(val, true); err != nil {
				return err
			}
		case Aggregate:
			if err := s.buildAggregate(val, true); err != nil {
				return err
			}
		case RawExpr:
			s.raw(val)
		default:
			return errs.NewErrUnsupportedSelectable(c)
		}
	}
	return nil
}

func (s *Selector[T]) buildColumn(c Column, useAlias bool) error {
	err := s.builder.buildColumn(c.name)
	if err != nil {
		return err
	}
	if useAlias {
		s.buildAs(c.alias)
	}
	return nil
}

// Where 用于构造 WHERE 查询条件。如果 ps 长度为 0，那么不会构造 WHERE 部分
func (s *Selector[T]) Where(ps ...Predicate) *Selector[T] {
	s.where = ps
	return s
}

// GroupBy 设置 group by 子句
func (s *Selector[T]) GroupBy(cols ...Column) *Selector[T] {
	s.groupBy = cols
	return s
}

func (s *Selector[T]) Having(ps ...Predicate) *Selector[T] {
	s.having = ps
	return s
}

func (s *Selector[T]) Offset(offset int) *Selector[T] {
	s.offset = offset
	return s
}

func (s *Selector[T]) Limit(limit int) *Selector[T] {
	s.limit = limit
	return s
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	//q, err := s.Build()
	//if err != nil {
	//	return nil, err
	//}
	//// s.db 是我们定义的 DB
	//// s.db.db 则是 sql.DB
	//// 使用 QueryContext，从而和 GetMulti 能够复用处理结果集的代码
	//rows, err := s.db.queryContext(ctx, q.SQL, q.Args...)
	//if err != nil {
	//	return nil, err
	//}
	//
	//if !rows.Next() {
	//	return nil, ErrNoRows
	//}
	//
	//tp := new(T)
	//meta, err := s.db.r.Get(tp)
	//if err != nil {
	//	return nil, err
	//}
	//val := s.db.valCreator(tp, meta)
	//err = val.SetColumns(rows)
	//return tp, err

	handler := s.getHandler
	// 套上 midddleware
	mdls := s.db.mdls
	for i:=len(mdls)-1;i>=0;i-- {
		handler = mdls[i](handler)
	}
	qc := &QueryContext{
		Builder: s,
		Type: "SELECT",
	}
	res := handler(ctx, qc)
	if res.Result != nil {
		return res.Result.(*T), res.Err
	}
	return nil, res.Err
}

var _ HandleFunc = (&Selector[any]{}).getHandler

// 将 get 调用作为 middleware 最后一级调用
func (s *Selector[T]) getHandler(ctx context.Context, qc *QueryContext) *QueryResult {
	q, err := qc.Builder.Build()
	if err != nil {
		return &QueryResult{
			Err: err,
		}
	}
	// s.db 是我们定义的 DB
	// s.db.db 则是 sql.DB
	// 使用 QueryContext，从而和 GetMulti 能够复用处理结果集的代码
	rows, err := s.db.queryContext(ctx, q.SQL, q.Args...)
	if err != nil {
		return &QueryResult{
			Err: err,
		}
	}

	if !rows.Next() {
		return &QueryResult{
			Err: ErrNoRows,
		}
	}

	tp := new(T)
	meta, err := s.db.r.Get(tp)
	if err != nil {
		return &QueryResult{
			Err: err,
		}
	}
	val := s.db.valCreator(tp, meta)
	err = val.SetColumns(rows)
	return &QueryResult{
		Result: tp,
		Err: err,
	}
}


func (s *Selector[T]) GetMulti(ctx context.Context) (res []*T, err error) {
	//var db sql.DB
	//q, err := s.Build()
	//if err != nil {
	//	return nil, err
	//}
	//rows, err := db.QueryContext(ctx, q.SQL, q.Args...)
	//if err != nil {
	//	return nil, err
	//}
	//
	//tp := new(T)
	//meta, err := s.db.r.Get(tp)
	//if err != nil {
	//	return nil, err
	//}
	//val := s.db.valCreator(tp, meta)
	//for rows.Next() {
	//	tp = new(T)
	//	err = val.SetColumns(rows)
	//	if err != nil {
	//		return nil, err
	//	}
	//	res = append(res, tp)
	//}
	//return
	
	handler := s.getMultiHandler
	mdls := s.db.mdls
	for i:=len(mdls)-1;i>=0;i-- {
		handler = mdls[i](handler)
	}
	qc := &QueryContext{
		Builder: s,
		Type: "SELECT",
	}
	qr := handler(ctx, qc)
	
	return qr.Result.([]*T), qr.Err
}

var _ HandleFunc = (&Selector[any]{}).getMultiHandler

// 将 GetMulti 调用作为 middleware 最后一级调用
func (s *Selector[T]) getMultiHandler(ctx context.Context, qc *QueryContext) *QueryResult {
	var db sql.DB
	var res []*T
	q, err := s.Build()
	if err != nil {
		return &QueryResult{
			Err: err,
		}
	}
	rows, err := db.QueryContext(ctx, q.SQL, q.Args...)
	if err != nil {
		return &QueryResult{
			Err: err,
		}
	}

	tp := new(T)
	meta, err := s.db.r.Get(tp)
	if err != nil {
		return &QueryResult{
			Err: err,
		}
	}
	val := s.db.valCreator(tp, meta)
	for rows.Next() {
		tp = new(T)
		err = val.SetColumns(rows)
		if err != nil {
			return &QueryResult{
				Err: err,
			}
		}
		res = append(res, tp)
	}
	return &QueryResult{
		Result: res,
		Err: err,
	}
}



// NewSelector 构造 Selector
func NewSelector[T any](db *DB) *Selector[T] {
	return &Selector[T]{
		builder: builder{
			dialect: db.dialect,
			quoter:  db.dialect.quoter(),
		},
		db: db,
	}
}

// Selectable select 语句中 colume, rawexpr, aggregate 的抽象
type Selectable interface {
	selectable()
}
