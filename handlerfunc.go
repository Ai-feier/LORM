package LORM

import (
	"LORM/model"
	"context"
)

// Middleware 函数式中间件
type Middleware func(next HandleFunc) HandleFunc

type HandleFunc func(ctx context.Context, qc *QueryContext) *QueryResult

type QueryContext struct {
	// Type 声明查询类型。即 SELECT, UPDATE, DELETE 和 INSERT
	Type    string
	Builder QueryBuilder
	Model   *model.Model
}

type QueryResult struct {
	Result any
	Err    error
}
