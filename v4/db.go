//go:build v4
package orm

type DBOption func(*DB)

type DB struct {
	r *registry
}

func NewDB(opts...DBOption) (*DB, error) {
	db := &DB{
		r: &registry{},
	}
	for _, opt := range opts {
		opt(db)
	}
	return db, nil
}

func MustDB(opts...DBOption) *DB {
	db, err := NewDB(opts...)
	if err != nil {
		panic(err)
	}
	return db
}