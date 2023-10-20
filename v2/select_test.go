//go:build v2
package v2


import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelector_Build(t *testing.T) {
	testCases := []struct{
		name string
		q QueryBuilder
		wantQuery *Query
		wantErr error
	} {
		{
			name: "no from",
			q: NewSelector[TestModel](),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel`;",
			},
		},
		{
			name: "from",
			q: NewSelector[TestModel]().From("`test_model_t`"),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model_t`;",
			},
		},
		{
			// 调用 FROM，但是传入空字符串
			name: "empty from",
			q:    NewSelector[TestModel]().From(""),
			wantQuery: &Query{
				SQL: "SELECT * FROM `TestModel`;",
			},
		},
		{
			// 调用 FROM，同时出入看了 DB
			name: "with db",
			q:    NewSelector[TestModel]().From("`test_db`.`test_model`"),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_db`.`test_model`;",
			},
		},
		{
			// 单一简单条件
			name: "single and simple predicate",
			q: NewSelector[TestModel]().From("`test_model_t`"). 
				Where(C("Id").EQ(1)),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model_t` WHERE `Id` = ?;",
				Args: []any{1},
			},
		},
		{
			// 多个 predicate
			name: "multiple predicates",
			q: NewSelector[TestModel]().
				Where(C("Age").GT(18), C("Age").LT(35)),
				wantQuery: &Query{
					SQL:  "SELECT * FROM `TestModel` WHERE (`Age` > ?) AND (`Age` < ?);",
					Args: []any{18,35},
				},
		},
		{
			// 使用 AND
			name: "and",
			q: NewSelector[TestModel]().
				Where(C("Age").GT(18).And(C("Age").LT(35))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel` WHERE (`Age` > ?) AND (`Age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			// 使用 OR
			name: "or",
			q:    NewSelector[TestModel]().
				Where(C("Age").GT(18).Or(C("Age").LT(35))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel` WHERE (`Age` > ?) OR (`Age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			// not 
			name: "not",
			q: NewSelector[TestModel]().Where(Not(C("Age").GT(18))),
			wantQuery: &Query{
				SQL: "SELECT * FROM `TestModel` WHERE  NOT (`Age` > ?);",
				Args: []any{18},
			},
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, err := tc.q.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, query)
		})
	}
}