//go:build v1
package v1

import "database/sql"




















type TestModel struct {
	Id int64 
	FirseName string
	Age int8 
	LastName *sql.NullString
}