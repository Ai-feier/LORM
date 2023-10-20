//go:build v8
package orm

import (
	"LORM/v8/internal/valuer"
	"LORM/v8/model"
	"database/sql"
	"fmt"
)

type DBOption func(*DB)

type DB struct {
	r model.Registry
	db *sql.DB
	valCreator valuer.Creator
}

func Open(driver string, dsn string, opts ...DBOption) (*DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return OpenDB(db, opts...)
}

func OpenDB(db *sql.DB, opts...DBOption) (*DB, error) {
	res := &DB {
		db: db,
		r: model.NewRegistry(),
		valCreator: valuer.NewUnsafeValue,
	}
	
	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}

func DBWithRegistry(r model.Registry) DBOption {
	return func(db *DB) {
		db.r = r
	}
}

// MustNewDB 创建一个 DB，如果失败则会 panic
func MustNewDB(driver string, dsn string, opts ...DBOption) *DB {
	db, err := Open(driver, dsn, opts...)
	if err != nil {
		panic(err)
	}
	return db
}


func getWordsInLongestSubsequence(n int, words []string, groups []int) []string {
	if len(words) == 1 {
		return []string{words[0]}
	}
	gn := len(groups)
	dp := make([]int, gn)

	ans := []int{}
	for i:=1;i<gn;i++ {
		tmp := []int{}
		for j := 0; j <= i; j++ {
			if len(words[i])==len(words[j]) && (len(tmp)==0 || groups[tmp[len(tmp)-1]] != groups[j]) {
				cnt := 0
				for k := 0; k < len(words[i]); k++ {
					if words[i][k] != words[j][k] {
						cnt++
					}
				}
				if cnt != 1 {
					continue
				}
				tmp = append(tmp, j)
				if len(tmp) > len(ans) {
					ans = tmp
				}
				dp[i] = dp[j]+1
			}
		}
	}
	fmt.Println(dp)

	var res []string
	for _, a := range ans {
		res = append(res, words[a])
	}
	if len(res) == 0 {
		res = append(res,words[0])
	}
	return res
}
// 找到最长不相等子序列的对应的下标



//3
//["e","a","b"]
//[0,0,1]
//4
//["a","b","c","d"]
//[1,0,1,1]
//1
//["c"]
//[0]
//2
//["d","g"]
//[0,1]

func findIndices(nums []int, in int, va int) []int {
	if in > len(nums) {
		return []int{-1,-1}
	}

	n := len(nums)
	pre := make([][2]int, n)
	pre[0][0] = nums[0]
	for i:=1;i<n;i++ {
		if nums[i] > pre[i-1][0] {
			pre[i][0] = nums[i]
			pre[i][1] = i
		} else {
			pre[i][0] = pre[i-1][0]
			pre[i][1] = pre[i-1][1]
		}
	}

	suf := make([][2]int, n)
	suf[n-1][0] = nums[len(nums)-1]
	suf[n-1][1] = n-1
	for i:=n-2;i>=0;i-- {
		if nums[i] > suf[i+1][0] {
			suf[i][0] = nums[i]
			suf[i][1] = i
		} else {
			suf[i][0] = suf[i+1][0]
			suf[i][1] = suf[i+1][1]
		}
	}
	fmt.Println("-----------",pre,suf)
	for i:=0;i<n;i++ {
		if abs(nums[i]-pre[i][0]) >= va && i-pre[i][1] >=in {
			return []int{pre[i][1], i}
		}
		if abs(nums[i]-suf[i][0]) >=va && suf[i][0]-i >=in {
			return []int{i, suf[i][1]}
		}
		
		//if abs(pre[i][0]-suf[i][0]) >= va {
		//	if abs(pre[i][1]-suf[i][1]) >= in {
		//		return []int{pre[i][1], suf[i][1]}
		//	}
		//}
	}
	return []int{-1,-1}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}