package service

import (
	"context"
	"fmt"
	"github.com/Ai-feier/rbacapp/dao"
	"github.com/Ai-feier/rbacapp/model"
)

type RoleBindingSvc struct {
}

func NewRoleBindingSvc() *RoleBindingSvc {
	return &RoleBindingSvc{}
}

func (r *RoleBindingSvc) CreateRoleBinding(ctx context.Context, rb *model.RoleBindings, role *model.Roles) (err error){
	if role.Namespace != rb.Namespace {
		return fmt.Errorf("不属于同一 namespace 将不建立 namespace")
	}

	rbDao := dao.NewRoleBindingDao(ctx)
	rb.RoleID = role.ID

	// TODO: 验证 role 是否存在

	return rbDao.CreateRoleBinding(rb)
}

func (r *RoleBindingSvc) DeleteRoleBinding(ctx context.Context, rb *model.RoleBindings) error {
	rbDao := dao.NewRoleBindingDao(ctx)
	return rbDao.DeleteRoleBinding(rb)
}
