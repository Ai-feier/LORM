package service

import (
	"context"
	"fmt"
	"github.com/Ai-feier/rbacapp/dao"
	"github.com/Ai-feier/rbacapp/model"
)

// 新建 role, 需同时建 rolesubref
// 删除 role, 同时删除 rolesubref

type RoleSvc struct {
}

func NewRoleSvc() *RoleSvc {
	return &RoleSvc{}
}

func (r *RoleSvc) CreateRole(ctx context.Context, role *model.Roles, subs ...*model.RoleSubRefs) error {
	roleDao := dao.NewRoleDao(ctx)
	err := roleDao.CreateRole(role)
	if err != nil {
		return err
	}

	// 同时新建 rolesubref
	subRefDao := dao.NewRoleSubRefDao(ctx)
	for _, sub := range subs {
		sub.RoleID = role.ID
		er := subRefDao.CreateRoleSubRef(sub)
		if er != nil {
			// 同时删除 role
			roleDao.DeleteRole(role)
			return fmt.Errorf("role: %s 创建失败: %w", role.Name, er)
		}
	}
	return nil
}

func (r *RoleSvc) DeleteRole(ctx context.Context, role *model.Roles) error {
	roleDao := dao.NewRoleDao(ctx)
	role, err := roleDao.FindByNameAndNameSpace(role.Name, role.Namespace)
	if role == nil {
		return fmt.Errorf("role 不存在")
	}
	if err != nil {
		return err
	}

	// 同时删除 subref
	subRefDao := dao.NewRoleSubRefDao(ctx)
	err = subRefDao.DeleteByRoleID(role.ID)
	if err != nil {
		return err
	}
	// 在删除 role
	return roleDao.DeleteRole(role)
}
