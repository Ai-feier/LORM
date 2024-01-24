package lorm

import (
	"context"
	"github.com/Ai-feier/lorm/model"
)

// Middleware 函数式中间件
type Middleware func(next HandleFunc) HandleFunc

type HandleFunc func(ctx context.Context, qc *QueryContext) *QueryResult

// QueryContext Middleware 向下传递的 Context
type QueryContext struct {
	// Type 声明查询类型。即 SELECT, UPDATE, DELETE 和 INSERT
	Type    string
	Builder QueryBuilder
	// Model 向 Middleware 提供当前 sql 操作的元数据信息
	Model   *model.Model
}

// QueryResult 每个 Handler 向上返回的结果
type QueryResult struct {
	Result any
	Err    error
}
