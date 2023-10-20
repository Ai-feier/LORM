//go:build v2
package v2

import "database/sql"

type TestModel struct {
	Id int64 
	FirstName string
	Age int8 
	LastName *sql.NullString
}