package dao

import (
	"context"
	"github.com/Ai-feier/rbacapp/model"
	"gorm.io/gorm"
)

type ClusterRoleSubRefDao struct {
	db *gorm.DB
}

func NewClusterRoleSubRefDao(ctx context.Context) *ClusterRoleSubRefDao {
	return &ClusterRoleSubRefDao{db: GetDB(ctx)}
}

func (c *ClusterRoleSubRefDao) CreateClusterRoleSubRef(sub *model.ClusterRoleSubRef) error {
	return c.db.Create(&sub).Error
}

func (c *ClusterRoleSubRefDao) DeleteClusterRoleSubRef(sub *model.ClusterRoleSubRef) error {
	return c.db.Model(&model.ClusterRoleSubRef{}).Where("id = ?", sub.ID).Error
}

func (c *ClusterRoleSubRefDao) DeleteByClusterID(clusterRoleID int) error {
	return c.db.Model(&model.ClusterRoleSubRef{}).Where("cluster_role_id = ?", clusterRoleID).
		Delete(&model.ClusterRoleSubRef{}).Error
}
