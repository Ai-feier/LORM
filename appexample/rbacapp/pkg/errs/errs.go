package errs

import "fmt"

func ExistUserName(name string) error {
	return fmt.Errorf("用户名: %s 已存在", name)
}
