package dao

import (
	"context"
	"github.com/Ai-feier/rbacapp/model"
	"gorm.io/gorm"
)

type RoleSubRefDao struct {
	db *gorm.DB
}

func NewRoleSubRefDao(ctx context.Context) *RoleSubRefDao {
	return &RoleSubRefDao{db: GetDB(ctx)}
}

func (r *RoleSubRefDao) CreateRoleSubRef(sub *model.RoleSubRef) error {
	return r.db.Model(&model.RoleSubRef{}).Create(&sub).Error
}

func (r *RoleSubRefDao) DeleteRoleSubRef(sub *model.RoleSubRef) error {
	return r.db.Delete(&sub).Error
}

func (r *RoleSubRefDao) DeleteByRoleID(roleID int) error  {
	return r.db.Model(&model.RoleSubRef{}).Where("role_id = ?", roleID).
		Delete(&model.RoleSubRef{}).Error
}
