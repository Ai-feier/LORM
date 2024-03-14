package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/Ai-feier/rbacapp/model"
	"gorm.io/gorm"
)

type RoleDao struct {
	db *gorm.DB
}

func NewRoleDao(ctx context.Context) *RoleDao {
	return &RoleDao{db: GetDB(ctx)}
}

func (r *RoleDao) FindByNameAndNameSpace(name, namespace string) (role *model.Role, err error) {
	var cnt int64
	err = r.db.Model(&model.Role{}).Where("name = ? and namespace = ?", name, namespace).Find(&role).Count(&cnt).Error
	if cnt == 0 {
		return nil, errors.New(fmt.Sprintf("当前 %s 下不存在 %s role", namespace, name))
	}
	return
}

func (r *RoleDao) CreateRole(role *model.Role) error {
	// 判断同一 namespace 下是否重名
	var cnt int64
	r.db.Model(&model.Role{}).Where("name = ?", role.Name).Count(&cnt)
	if cnt > 0 {
		return fmt.Errorf("当前命名空间 %s 下已存在名为 %s 的 role, 请重命名", role.Namespace, role.Name)
	}
	return r.db.Model(&model.Role{}).Create(&role).Error
}

func (r *RoleDao) DeleteRole(role *model.Role) error {
	return r.db.Model(&model.Role{}).Where("name = ? and namespace = ?", role.Name, role.Namespace).
		Delete(&model.Role{}).Error

	// 删除前检查是否存在
	//var cnt int64
	//r.db.Model(&model.Role{}).Where("name = ? and namespace = ?", role.Name, role.Namespace).Count(&cnt)
	//if cnt == 0 {
	//	return fmt.Errorf("前命名空间 %s 下不存在名为 %s 的 role", role.Namespace, role.Name)
	//}
	//return r.db.Model(&model.Role{}).Where("name = ? and namespace = ?", role.Name, role.Namespace).Delete(&model.Role{}).Error
}

func (r *RoleDao) DeleteRoleByNameAndNamespace(name, namespace string) error {
	return r.db.Model(&model.Role{}).Where("name = ? and namespace = ?", name, namespace).
		Delete(&model.Role{}).Error
}
