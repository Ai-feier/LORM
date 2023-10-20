package sql_demo

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JsonColumn[T any] struct {
	Val T
	
	// 解决 NULL
	Valid bool 
}

func (j JsonColumn[T]) Value() (driver.Value, error) {
	// 没有数据
	if !j.Valid {
		return nil, nil
	}
	return json.Marshal(j.Val)
}

func (j *JsonColumn[T]) Scan(src any) error {
	var bs []byte
	switch data := src.(type) {
	case []byte:
		bs = data
	case string:
		bs = []byte(data)
	case nil:
		return nil
	default:
		return fmt.Errorf("不支持的数据类型 %v", src)
	}
	err := json.Unmarshal(bs, &j.Val)
	if err == nil {
		j.Valid = true
	}
	return err
}



