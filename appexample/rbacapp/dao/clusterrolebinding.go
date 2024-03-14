package dao

import (
	"context"
	"fmt"
	"github.com/Ai-feier/rbacapp/model"
	"gorm.io/gorm"
)

type ClusterRoleBindingDao struct {
	db *gorm.DB
}

func NewClusterRoleBindingDao(ctx context.Context) *ClusterRoleBindingDao {
	return &ClusterRoleBindingDao{db: GetDB(ctx)}
}

func (c *ClusterRoleBindingDao) CreateClusterRoleBinding(b *model.ClusterRoleBinding) error {
	var cnt int64
	c.db.Model(&model.ClusterRoleBinding{}).Where("name = ?", b.Name).Count(&cnt)
	if cnt > 0 {
		return fmt.Errorf("已存在名为 %s 的 clusterrolebinding, 请重命名",b.Name)
	}
	return c.db.Create(b).Error
}

func (c *ClusterRoleBindingDao) DeleteClusterRoleBinding(b *model.ClusterRoleBinding) error {
	return c.db.Delete(&b).Error
	//return c.db.Model(&model.ClusterRole{}).Where("id = ?", b.ID).Delete(b).Error
}

func (c *ClusterRoleBindingDao) DeleteClusterRoleBindingByName(name string) error {
	return c.db.Model(&model.ClusterRoleBinding{}).Where("name = ?", name).Delete(&model.ClusterRoleBinding{}).Error
}

func (c *ClusterRoleBindingDao) FindByName(name string) (crb *model.ClusterRoleBinding, err error) {
	err = c.db.Model(&model.ClusterRoleBinding{}).Where("name = ?", name).Find(&crb).Error
	return
}
