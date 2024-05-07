package dao

import (
	"context"
	"github.com/Ai-feier/lorm"
	"github.com/Ai-feier/rbacapp/model"
)

type ClusterRoleDao struct {
	db *lorm.DB
}

func NewClusterRoleDao(ctx context.Context) *ClusterRoleDao {
	return &ClusterRoleDao{db: GetDB(ctx)}
}

func (c *ClusterRoleDao) CreateClusterRole(role *model.ClusterRoles) error {
	// 是否存在
	//var (
	//	cnt int64
	//)
	////c.db.Model(&model.ClusterRoles{}).Where("name = ?", role.Name).Count(&cnt)
	//if cnt > 0 {
	//	return fmt.Errorf("已存在名为 %s 的 clusterrole, 请重命名", role.Name)
	//}
	//return c.db.Create(role).Error

	res := lorm.NewInserter[model.ClusterRoles](c.db).Values(role).Exec(context.Background())
	return res.Err()
}

func (c *ClusterRoleDao) DeleteClusterRole(role *model.ClusterRoles) error {
	//return c.db.Model(&model.ClusterRoles{}).Where("id = ?", role.ID).Delete(&role).Error

	//return c.db.Model(&model.ClusterRoles{}).Where("name = ?", role.Name).Delete(role).Error

	// lorm
	res := lorm.NewDeleter[model.ClusterRoles](c.db).Where(lorm.C("ID").EQ(role.ID)).Exec(context.Background())
	return res.Err()
}

func (c *ClusterRoleDao) FindByName(name string) (cr *model.ClusterRoles, err error) {
	//var cnt int64
	//err = c.db.Model(&model.ClusterRoles{}).Where("name = ?", name).Find(&cr).Count(&cnt).Error
	//if cnt == 0 {
	//	return nil, errors.New("clusterrole 不存在")
	//}
	//return

	res, err := lorm.NewSelector[model.ClusterRoles](c.db).Where(lorm.C("Name").EQ(name)).Get(context.Background())
	return res, err
}
