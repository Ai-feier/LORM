package lorm

import (
	"context"
	"github.com/Ai-feier/lorm/internal/errs"
	"github.com/Ai-feier/lorm/model"
)

type UpsertBuilder[T any] struct {
	i               *Inserter[T]
	conflictColumns []string
}

type Upsert struct {
	conflictColumns []string
	assigns         []Assignable
}

func (o *UpsertBuilder[T]) ConflictColumns(cols ...string) *UpsertBuilder[T] {
	o.conflictColumns = cols
	return o
}

// Update 也可以看做是一个终结方法，重新回到 Inserter 里面
func (o *UpsertBuilder[T]) Update(assigns ...Assignable) *Inserter[T] {
	o.i.upsert = &Upsert{
		conflictColumns: o.conflictColumns,
		assigns:         assigns,
	}
	return o.i
}

type Inserter[T any] struct {
	builder
	values  []*T
	db      *DB
	columns []string
	upsert  *Upsert
}

func (i *Inserter[T]) Exec(ctx context.Context) Result {
	//q, err := i.Build()
	//if err != nil {
	//	return Result{err: err}
	//}
	//res, err := i.db.execContext(ctx, q.SQL, q.Args...)
	//return Result{res: res, err: err}


	handler := i.execHandler
	mdls := i.db.mdls
	for j:=len(mdls)-1;j>=0;j-- {
		handler = mdls[j](handler)
	}
	qc := &QueryContext{
		Builder: i,
		Type: "INSERT",
		Model: i.model,
	}
	qr := handler(ctx, qc)

	return qr.Result.(Result)
}

func (i *Inserter[T]) execHandler(ctx context.Context, qc *QueryContext) *QueryResult {
	q, err := i.Build()
	if err != nil {
		return &QueryResult{
			Result: Result{err: err},
		}
	}
	res, err := i.db.execContext(ctx, q.SQL, q.Args...)
	return &QueryResult{
		Result: Result{res: res, err: err},
	}
}

func NewInserter[T any](db *DB) *Inserter[T] {
	return &Inserter[T]{
		db: db,
		builder: builder{
			dialect: db.dialect,
			quoter:  db.dialect.quoter(),
		},
	}
}

func (i *Inserter[T]) Values(vals ...*T) *Inserter[T] {
	i.values = vals
	return i
}

func (i *Inserter[T]) OnDuplicateKey() *UpsertBuilder[T] {
	return &UpsertBuilder[T]{
		i: i,
	}
}

// Fields 指定要插入的列
func (i *Inserter[T]) Columns(cols ...string) *Inserter[T] {
	i.columns = cols
	return i
}

func (i *Inserter[T]) Build() (*Query, error) {
	if len(i.values) == 0 {
		return nil, errs.ErrInsertZeroRow
	}
	m, err := i.db.r.Get(i.values[0])
	if err != nil {
		return nil, err
	}
	i.model = m

	i.sb.WriteString("INSERT INTO ")
	i.quote(m.TableName)
	i.sb.WriteString("(")

	// 插入字段
	fields := m.Fields
	if len(i.columns) != 0 {
		fields = make([]*model.Field, 0, len(i.columns))
		for _, c := range i.columns {
			field, ok := m.FieldMap[c]
			if !ok {
				return nil, errs.NewErrUnknownField(c)
			}
			fields = append(fields, field)
		}
	}

	// (len(i.values) + 1) 中 +1 是考虑到 UPSERT 语句会传递额外的参数
	i.args = make([]any, 0, len(fields)*(len(i.values)+1))
	for idx, fd := range fields {
		if idx > 0 {
			i.sb.WriteByte(',')
		}
		i.quote(fd.ColName)
	}

	i.sb.WriteString(") VALUES")
	for vIdx, val := range i.values {
		if vIdx > 0 {
			i.sb.WriteByte(',')
		}
		refVal := i.db.valCreator(val, i.model)
		i.sb.WriteByte('(')
		for fIdx, field := range fields {
			if fIdx > 0 {
				i.sb.WriteByte(',')
			}
			i.sb.WriteByte('?')
			fdVal, err := refVal.Field(field.GoName)
			if err != nil {
				return nil, err
			}
			i.addArgs(fdVal)
		}
		i.sb.WriteByte(')')
	}

	if i.upsert != nil {
		err = i.dialect.buildUpsert(&i.builder, i.upsert)
		if err != nil {
			return nil, err
		}
	}

	i.sb.WriteString(";")
	return &Query{
		SQL:  i.sb.String(),
		Args: i.args,
	}, nil
}
