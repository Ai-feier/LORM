package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/Ai-feier/rbacapp/model"
	"gorm.io/gorm"
)

// 传入 role_id, user_id 创建 rolebinding
// 删除 rolebinding

type RoleBindingDao struct {
	db *gorm.DB
}

func NewRoleBindingDao(ctx context.Context) *RoleBindingDao {
	return &RoleBindingDao{db: GetDB(ctx)}
}

// 构建好了传进来
func (r *RoleBindingDao) CreateRoleBinding(b *model.RoleBinding) error {
	var cnt int64
	r.db.Model(&model.RoleBinding{}).Where("namespace = ? and name = ?", b.Namespace, b.Name).Count(&cnt)
	if cnt > 0 {
		return fmt.Errorf("当前命名空间 %s 下已存在名为 %s 的 rolebinding, 请检查配置清单", b.Namespace, b.Name)
	}
	return r.db.Model(&model.RoleBinding{}).Create(&b).Error
}

func (r *RoleBindingDao) DeleteRoleBinding(b *model.RoleBinding) error {
	var cnt int64
	err :=r.db.Model(&model.RoleBinding{}).Where("namespace = ? and name = ?", b.Namespace, b.Name).Count(&cnt).
		Delete(&model.RoleBinding{}).Error
	if cnt == 0 {
		return errors.New("当前 rolebinding 不存在")
	}
	return err
}

func (r RoleBindingDao) FindByNameAndNamespace(name, namespace string) (rb *model.RoleBinding, err error){
	err = r.db.Model(&model.RoleBinding{}).Where("name = ? and namespace = ?", name, namespace).Find(&rb).Error
	return
}
