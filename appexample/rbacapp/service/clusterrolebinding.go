package service

import (
	"context"
	"github.com/Ai-feier/rbacapp/dao"
	"github.com/Ai-feier/rbacapp/model"
)

type ClusterRoleBindingSvc struct {
}

func NewClusterRoleBindingSvc() *ClusterRoleBindingSvc {
	return &ClusterRoleBindingSvc{}
}

func (c *ClusterRoleBindingSvc) CreateClusterRoleBinding(ctx context.Context, name, users string, cr *model.ClusterRole) (err error){
	crb := &model.ClusterRoleBinding{
		RoleID: cr.ID,
		Users: users,
		Name: name,
	}

	// TODO: 验证 role 是否存在

	return dao.NewClusterRoleBindingDao(ctx).CreateClusterRoleBinding(crb)
}

func (c *ClusterRoleBindingSvc) DeleteClusterRoleBindingByName(ctx context.Context, name string) (err error) {
	crbDao := dao.NewClusterRoleBindingDao(ctx)
	err = crbDao.DeleteClusterRoleBindingByName(name)
	return err
}
