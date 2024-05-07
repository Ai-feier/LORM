package dao

import (
	"context"
	"github.com/Ai-feier/lorm"
	"github.com/Ai-feier/rbacapp/model"
)

type RoleSubRefDao struct {
	db *lorm.DB
}

func NewRoleSubRefDao(ctx context.Context) *RoleSubRefDao {
	return &RoleSubRefDao{db: GetDB(ctx)}
}

func (r *RoleSubRefDao) CreateRoleSubRef(sub *model.RoleSubRefs) error {
	//return r.db.Model(&model.RoleSubRefs{}).Create(&sub).Error

	return lorm.NewInserter[model.RoleSubRefs](r.db).Values(sub).Exec(context.Background()).Err()
}

func (r *RoleSubRefDao) DeleteRoleSubRef(sub *model.RoleSubRefs) error {
	//return r.db.Delete(&sub).Error

	return lorm.NewDeleter[model.RoleSubRefs](r.db).Where(lorm.C("ID").EQ(sub.ID)).Exec(context.Background()).Err()
}

func (r *RoleSubRefDao) DeleteByRoleID(roleID int) error  {
	//return r.db.Model(&model.RoleSubRefs{}).Where("role_id = ?", roleID).
	//	Delete(&model.RoleSubRefs{}).Error

	return lorm.NewDeleter[model.RoleSubRefs](r.db).Where(lorm.C("ID").EQ(roleID)).Exec(context.Background()).Err()
}
