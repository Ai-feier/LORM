package reflect

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIterateArray(t *testing.T) {
	testCases := []struct{
		name string
		entity any 
		
		wantRes []any 
		wantErr error
	} {
		{
			name: "array",
			entity: [3]int{1, 2, 3},
			wantRes: []any{1, 2, 3},
		},
		{
			name: "slice",
			entity: []int{1, 2, 3},
			wantRes: []any{1, 2, 3},
		},

	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err :=  IterateArrayOrSlice(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return 
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestIterateMap(t *testing.T) {
	testCases := []struct{
		name string 
		entity any 
		
		wantKeys []any 
		wantValues []any 
		wantErr error
	} {
		{
			name: "map",
			entity: map[string]string{
				"A":"a",
				"B":"b",
			},
			
			wantKeys: []any{"A","B"},
			wantValues: []any{"a", "b"},
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			key, val, err := IterateMap(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, len(tc.wantKeys), len(key))
			assert.Equal(t, len(tc.wantValues), len(val))
		})
	}
}