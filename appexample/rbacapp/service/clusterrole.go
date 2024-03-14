package service

import (
	"context"
	"fmt"
	"github.com/Ai-feier/rbacapp/dao"
	"github.com/Ai-feier/rbacapp/model"
)

type ClusterRoleSvc struct {
}

func NewClusterRoleSvc() *ClusterRoleSvc {
	return &ClusterRoleSvc{}
}

func (c *ClusterRoleSvc) CreateClusterRole(ctx context.Context, cr *model.ClusterRole, subs ...*model.ClusterRoleSubRef) (err error)  {
	crDao := dao.NewClusterRoleDao(ctx)
	//cr, err = crDao.FindByName(name)
	//// 已存在
	//if cr != nil {
	//	return fmt.Errorf("已存在名为 %s 的 clusterrole, 请重命名", name)
	//}

	// 创建 clusterrole
	err = crDao.CreateClusterRole(cr)
	if err != nil {
		return err
	}

	// 创建 subref
	subDao := dao.NewClusterRoleSubRefDao(ctx)
	for _, sub := range subs {
		sub.ClusterRoleID = cr.ID
		er := subDao.CreateClusterRoleSubRef(sub)
		if er != nil {
			// 删除 clusterrole
			crDao.DeleteClusterRole(cr)
			return fmt.Errorf("clusterrole: %s 创建失败: %w", cr.Name, er)
		}
	}
	return
}

func (c *ClusterRoleSvc) DeleteClusterRole(ctx context.Context, name string) error {
	crDao := dao.NewClusterRoleDao(ctx)
	cr, err := crDao.FindByName(name)
	if cr == nil {
		return fmt.Errorf("clusterrole 不存在")
	}
	if err != nil {
		return err
	}

	// 同时删除 subref
	subDao := dao.NewClusterRoleSubRefDao(ctx)
	err = subDao.DeleteByClusterID(cr.ID)
	if err != nil {
		return err
	}
	return crDao.DeleteClusterRole(cr)
}
