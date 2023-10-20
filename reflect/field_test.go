package reflect

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIterateFields(t *testing.T) {
	type User struct {
		Name string
		age int
	}
	testCases := []struct{
		name string 
		
		entity any 
		
		wantRes map[string]any 
		wantErr error
	} {
		{
			name: "struct",
			
			entity: User{
				Name: "Tom",
				age:  18,
			},
			wantRes: map[string]any{
				"Name": "Tom",
				"age": 0,
			},
		},
		{
			name: "point struct",
			entity: User{
				Name: "Tom",
				age:  18,
			},
			wantRes: map[string]any{
				"Name": "Tom",
				"age": 0,
			},

		},
		{
			name: "nil",
			entity: nil,
			wantErr: errors.New("不支持 nil"),
			
		},
		{
			name: "multi point",

			entity: func() **User {
				res := &User{
					Name: "Tom",
					age:  18,
				}
				return &res
			}(),
			wantRes: map[string]any{
				"Name": "Tom",
				"age": 0,
			},
		},
		{
			name: "zero nil",

			entity: (*User)(nil),
			wantErr: errors.New("不支持零值"),
		},	
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := IterateFields(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestSetField(t *testing.T) {
	type User struct {
		Name string
		age int
	}
	testCases := []struct{
		name string
		
		entity any 
		field string
		newValue any
		
		wantErr error
		wantEntity any
		
	} {
		{
			name: "struct",
			entity: User{
				Name: "Tom",
			},
			field: "Name",
			newValue: "Jerry",
			
			wantErr: errors.New("不可修改"),
		},
		{
			name: "pointer struct",
			entity: &User{
				Name: "Tom",
			},
			field: "Name",
			newValue: "Jerry",
			
			wantEntity: &User{
				Name: "Jerry",
			},

		},
		
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := SetField(tc.entity, tc.field, tc.newValue)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantEntity, tc.entity)
		})
	}
}