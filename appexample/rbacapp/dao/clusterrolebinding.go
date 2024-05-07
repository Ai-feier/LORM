package dao

import (
	"context"
	"fmt"
	"github.com/Ai-feier/lorm"
	"github.com/Ai-feier/rbacapp/model"
)

type ClusterRoleBindingDao struct {
	db *lorm.DB
}

func NewClusterRoleBindingDao(ctx context.Context) *ClusterRoleBindingDao {
	return &ClusterRoleBindingDao{db: GetDB(ctx)}
}

func (c *ClusterRoleBindingDao) CreateClusterRoleBinding(b *model.ClusterRoleBindings) error {
	//var cnt int64
	//c.db.Model(&model.ClusterRoleBindings{}).Where("name = ?", b.Name).Count(&cnt)
	//if cnt > 0 {
	//	return fmt.Errorf("已存在名为 %s 的 clusterrolebinding, 请重命名",b.Name)
	//}
	//return c.db.Create(b).Error

	res, _ := lorm.NewSelector[model.ClusterRoleBindings](c.db).Where(lorm.C("Name").EQ(b.Name)).Get(context.Background())
	if res == nil {
			return fmt.Errorf("已存在名为 %s 的 clusterrolebinding, 请重命名",b.Name)
	}
	ans := lorm.NewInserter[model.ClusterRoleBindings](c.db).Values(b).Exec(context.Background())
	return ans.Err()
}

func (c *ClusterRoleBindingDao) DeleteClusterRoleBinding(b *model.ClusterRoleBindings) error {
	//return c.db.Delete(&b).Error
	//return c.db.Model(&model.ClusterRoles{}).Where("id = ?", b.ID).Delete(b).Error

	res := lorm.NewDeleter[model.ClusterRoleBindings](c.db).Where(lorm.C("Name").EQ(b.ID)).Exec(context.Background())
	return res.Err()
}

func (c *ClusterRoleBindingDao) DeleteClusterRoleBindingByName(name string) error {
	//return c.db.Model(&model.ClusterRoleBindings{}).Where("name = ?", name).Delete(&model.ClusterRoleBindings{}).Error


	res := lorm.NewDeleter[model.ClusterRoleBindings](c.db).Where(lorm.C("Name").EQ(name)).Exec(context.Background())
	return res.Err()
}

func (c *ClusterRoleBindingDao) FindByName(name string) (crb *model.ClusterRoleBindings, err error) {
	//err = c.db.Model(&model.ClusterRoleBindings{}).Where("name = ?", name).Find(&crb).Error

	crb, err = lorm.NewSelector[model.ClusterRoleBindings](c.db).Where(lorm.C("Name").EQ(name)).Get(context.Background())
	return
}
