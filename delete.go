package lorm

import (
	"context"
)

type Deleter[T any] struct {
	builder
	tableName string
	where     []Predicate
	db        *DB
}

func NewDeleter[T any](db *DB) *Deleter[T] {
	return &Deleter[T]{
		db: db,
		builder: builder{
			quoter:  db.dialect.quoter(),
			dialect: db.dialect,
			r: db.r,
		},
	}
}

func (d *Deleter[T]) Build() (*Query, error) {
	defer func() {
		d.sb.Reset()
	}()
	var (
		t   T
		err error
	)
	d.model, err = d.db.r.Get(&t)
	if err != nil {
		return nil, err
	}
	d.sb.WriteString("DELETE FROM ")
	if d.tableName == "" {
		d.quote(d.model.TableName)
	} else {
		d.sb.WriteString(d.tableName)
	}
	if len(d.where) > 0 {
		d.sb.WriteString(" WHERE ")
		if err = d.buildPredicates(d.where); err != nil {
			return nil, err
		}
	}
	d.sb.WriteByte(';')
	return &Query{
		SQL:  d.sb.String(),
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

func (d *Deleter[T]) Exec(ctx context.Context) Result {
	//query, err := d.Build()
	//if err != nil {
	//	return Result{err: err}
	//}
	//res, err := d.db.execContext(ctx, query.SQL, d.args...)
	//return Result{res: res, err: err}

	handler := d.execHandler
	mdls := d.db.mdls
	for i:=len(mdls)-1;i>=0;i-- {
		handler = mdls[i](handler)
	}
	qc := &QueryContext{
		Builder: d,
		Type: "DELETE",
		Model: d.model,
	}
	qr := handler(ctx, qc)
	return qr.Result.(Result)
}

func (d *Deleter[T]) execHandler(ctx context.Context, qc *QueryContext) *QueryResult {
	query, err := d.Build()
	if err != nil {
		return &QueryResult{
			Result: Result{err: err},
		}
	}
	res, err := d.db.execContext(ctx, query.SQL, d.args...)
	return &QueryResult{
		Result: Result{res: res, err: err},
	}
}
