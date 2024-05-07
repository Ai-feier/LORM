package dao

import (
	"context"
	"github.com/Ai-feier/lorm"
	"github.com/Ai-feier/rbacapp/model"
)

type ClusterRoleSubRefDao struct {
	db *lorm.DB
}

func NewClusterRoleSubRefDao(ctx context.Context) *ClusterRoleSubRefDao {
	return &ClusterRoleSubRefDao{db: GetDB(ctx)}
}

func (c *ClusterRoleSubRefDao) CreateClusterRoleSubRef(sub *model.ClusterRoleSubRefs) error {
	//return c.db.Create(&sub).Error

	return lorm.NewInserter[model.ClusterRoleSubRefs](c.db).Values(sub).Exec(context.Background()).Err()
}

func (c *ClusterRoleSubRefDao) DeleteClusterRoleSubRef(sub *model.ClusterRoleSubRefs) error {
	//return c.db.Model(&model.ClusterRoleSubRefs{}).Where("id = ?", sub.ID).Error

	return lorm.NewDeleter[model.ClusterRoleSubRefs](c.db).Where(lorm.C("ID").EQ(sub.ID)).Exec(context.Background()).Err()
}

func (c *ClusterRoleSubRefDao) DeleteByClusterID(clusterRoleID int) error {
	//return c.db.Model(&model.ClusterRoleSubRefs{}).Where("cluster_role_id = ?", clusterRoleID).
	//	Delete(&model.ClusterRoleSubRefs{}).Error

	return lorm.NewDeleter[model.ClusterRoleSubRefs](c.db).Where(lorm.C("ClusterRoleID").EQ(clusterRoleID)).Exec(context.Background()).Err()
}
