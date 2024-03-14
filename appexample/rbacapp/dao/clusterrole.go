package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/Ai-feier/rbacapp/model"
	"gorm.io/gorm"
)

type ClusterRoleDao struct {
	db *gorm.DB
}

func NewClusterRoleDao(ctx context.Context) *ClusterRoleDao {
	return &ClusterRoleDao{db: GetDB(ctx)}
}

func (c *ClusterRoleDao) CreateClusterRole(role *model.ClusterRole) error {
	// 是否存在
	var (
		cnt int64
	)
	c.db.Model(&model.ClusterRole{}).Where("name = ?", role.Name).Count(&cnt)
	if cnt > 0 {
		return fmt.Errorf("已存在名为 %s 的 clusterrole, 请重命名", role.Name)
	}
	return c.db.Create(role).Error
}

func (c *ClusterRoleDao) DeleteClusterRole(role *model.ClusterRole) error {
	return c.db.Model(&model.ClusterRole{}).Where("id = ?", role.ID).Delete(&role).Error

	//return c.db.Model(&model.ClusterRole{}).Where("name = ?", role.Name).Delete(role).Error
}

func (c *ClusterRoleDao) FindByName(name string) (cr *model.ClusterRole, err error) {
	var cnt int64
	err = c.db.Model(&model.ClusterRole{}).Where("name = ?", name).Find(&cr).Count(&cnt).Error
	if cnt == 0 {
		return nil, errors.New("clusterrole 不存在")
	}
	return
}
