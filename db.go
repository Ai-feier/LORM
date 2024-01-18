package LORM

import (
	"LORM/internal/errs"
	"LORM/internal/valuer"
	"LORM/model"
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"log"
)

type DBOption func(*DB)

type DB struct {
	dialect    Dialect
	r          model.Registry
	db         *sql.DB
	valCreator valuer.Creator
	mdls       []Middleware
}

// Open 创建一个 DB 实例。
// 默认情况下，该 DB 将使用 MySQL 作为方言
// 如果你使用了其它数据库，可以使用 DBWithDialect 指定
func Open(driver string, dsn string, opts ...DBOption) (*DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return OpenDB(db, opts...)
}

func OpenDB(db *sql.DB, opts ...DBOption) (*DB, error) {
	res := &DB{
		dialect:    MySQL,
		r:          model.NewRegistry(),
		db:         db,
		valCreator: valuer.NewUnsafeValue,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}

func DBWithDialect(dialect Dialect) DBOption {
	return func(db *DB) {
		db.dialect = dialect
	}
}

// DBWithRegistry 更换数据中心实现
func DBWithRegistry(r model.Registry) DBOption {
	return func(db *DB) {
		db.r = r
	}
}

// DBUseReflectValuer 更改数据返回结果实现
func DBUseReflectValuer() DBOption {
	return func(db *DB) {
		db.valCreator = valuer.NewReflectValue
	}
}

// DBUseMiddlewares 注册 Middlewares
func DBUseMiddlewares(mdls ...Middleware) DBOption {
	return func(db *DB) {
		db.mdls = mdls
	}
}

// MustNewDB 创建一个 DB，如果失败则会 panic
func MustNewDB(driver string, dsn string, opts ...DBOption) *DB {
	db, err := Open(driver, dsn, opts...)
	if err != nil {
		panic(err)
	}
	return db
}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) Wait() error {
	err := db.db.Ping()
	for errors.Is(err, driver.ErrBadConn) {
		log.Println("数据库启动中")
		err = db.db.Ping()
	}
	return nil
}

func (db *DB) queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return db.db.QueryContext(ctx, query, args...)
}

func (db *DB) execContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return db.db.ExecContext(ctx, query, args...)
}

// BeginTx 开启事务
func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &Tx{tx: tx, db: db}, nil
}

// DoTx 将会开启事务执行 fn。如果 fn 返回错误或者发生 panic，事务将会回滚，
// 否则提交事务
func (db *DB) DoTx(ctx context.Context,
	fn func(ctx context.Context, tx *Tx) error,
	opts *sql.TxOptions) (err error) {
	var tx *Tx
	tx, err = db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	panicked := true
	defer func() {
		if panicked || err != nil {
			e := tx.Rollback()
			if e != nil {
				err = errs.NewErrFailToRollbackTx(err, e, panicked)
			}
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(ctx, tx)
	panicked = false
	return err
}
