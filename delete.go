package LORM

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
		},
	}
}

func (d *Deleter[T]) Build() (*Query, error) {
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
