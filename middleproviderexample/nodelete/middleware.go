package querylog

import (
	"context"
	"errors"
	orm "github.com/Ai-feier/lorm"
)

type MiddlewareBuilder struct {
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
	}
}

func (m MiddlewareBuilder) Build() orm.Middleware {
	return func(next orm.HandleFunc) orm.HandleFunc {
		return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
			// 禁用 DELETE 语句
			if qc.Type == "DELETE" {
				return &orm.QueryResult{
					Err: errors.New("禁止 Delete 语句"),
				}
			}
			return next(ctx, qc)
		}
	}
}
