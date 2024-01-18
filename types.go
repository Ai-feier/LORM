package lorm

import (
	"context"
)

// Querier sql 查询语句抽象
type Querier[T any] interface {
	Get(ctx context.Context) (*T, error)
	GetMulti(ctx context.Context) ([]*T, error)
}

// Executor insert update 语句执行抽象
type Executor interface {
	Exec(ctx context.Context) Result
}

// Query sql 中间结构体
// SQL: sql 语句
// Args: sql 语句中的占位符参数
type Query struct {
	SQL  string
	Args []any
}

// QueryBuilder 构造 sql 语句的抽象
type QueryBuilder interface {
	Build() (*Query, error)
}
