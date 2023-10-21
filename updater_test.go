package LORM

import (
	"LORM/internal/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdater_Build(t *testing.T) {
	db := memoryDB(t)
	testCases := []struct {
		name      string
		u         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			name:    "no column",
			u:       NewUpdater[TestModel](db),
			wantErr: errs.ErrNoUpdatedColumns,
		},
		{
			name: "single column",
			u: NewUpdater[TestModel](db).Update(&TestModel{
				Age: 18,
			}).Set(C("Age")),
			wantQuery: &Query{
				SQL:  "UPDATE `test_model` SET `age`=?;",
				Args: []any{int8(18)},
			},
		},
		{
			name: "assignment",
			u: NewUpdater[TestModel](db).Update(&TestModel{
				Age:       18,
				FirstName: "Tom",
			}).Set(C("Age"), Assign("FirstName", "YangZhuolin")),
			wantQuery: &Query{
				SQL:  "UPDATE `test_model` SET `age`=?,`first_name`=?;",
				Args: []any{int8(18), "YangZhuolin"},
			},
		},
		{
			name: "where",
			u: NewUpdater[TestModel](db).Update(&TestModel{
				Age:       18,
				FirstName: "Tom",
			}).Set(C("Age"), Assign("FirstName", "Yang")).
				Where(C("Id").EQ(1)),
			wantQuery: &Query{
				SQL:  "UPDATE `test_model` SET `age`=?,`first_name`=? WHERE `id` = ?;",
				Args: []any{int8(18), "Yang", 1},
			},
		},
		{
			name: "incremental",
			u: NewUpdater[TestModel](db).Update(&TestModel{
				Age:       18,
				FirstName: "Tom",
			}).Set(Assign("Age", C("Age").Add(1))),
			wantQuery: &Query{
				SQL:  "UPDATE `test_model` SET `age`=`age` + ?;",
				Args: []any{1},
			},
		},
		{
			name: "incremental-raw",
			u: NewUpdater[TestModel](db).Update(&TestModel{
				Age:       18,
				FirstName: "Tom",
			}).Set(Assign("Age", Raw("`age`+?", 1))),
			wantQuery: &Query{
				SQL:  "UPDATE `test_model` SET `age`=`age`+?;",
				Args: []any{1},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q, err := tc.u.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, q)
		})
	}
}
