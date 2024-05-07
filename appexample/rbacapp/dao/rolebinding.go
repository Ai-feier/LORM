package dao

import (
	"context"
	"fmt"
	"github.com/Ai-feier/lorm"
	"github.com/Ai-feier/rbacapp/model"
)

// 传入 role_id, user_id 创建 rolebinding
// 删除 rolebinding

type RoleBindingDao struct {
	db *lorm.DB
}

func NewRoleBindingDao(ctx context.Context) *RoleBindingDao {
	return &RoleBindingDao{db: GetDB(ctx)}
}

// 构建好了传进来
func (r *RoleBindingDao) CreateRoleBinding(b *model.RoleBindings) error {
	//var cnt int64
	//r.db.Model(&model.RoleBindings{}).Where("namespace = ? and name = ?", b.Namespace, b.Name).Count(&cnt)
	//if cnt > 0 {
	//	return fmt.Errorf("当前命名空间 %s 下已存在名为 %s 的 rolebinding, 请检查配置清单", b.Namespace, b.Name)
	//}
	//return r.db.Model(&model.RoleBindings{}).Create(&b).Error

	res, _ := lorm.NewSelector[model.RoleBindings](r.db).Where(lorm.C("Namespace").EQ(b.Namespace), lorm.C("Name").EQ(b.Name)).Get(context.Background())
	if res != nil {
			return fmt.Errorf("当前命名空间 %s 下已存在名为 %s 的 rolebinding, 请检查配置清单", b.Namespace, b.Name)
	}

	ans := lorm.NewInserter[model.RoleBindings](r.db).Values(b).Exec(context.Background())
	return ans.Err()
}

func (r *RoleBindingDao) DeleteRoleBinding(b *model.RoleBindings) error {
	//var cnt int64
	//err :=r.db.Model(&model.RoleBindings{}).Where("namespace = ? and name = ?", b.Namespace, b.Name).Count(&cnt).
	//	Delete(&model.RoleBindings{}).Error
	//if cnt == 0 {
	//	return errors.New("当前 rolebinding 不存在")
	//}
	//return err

	return lorm.NewDeleter[model.RoleBindings](r.db).Where(lorm.C("Namespace").EQ(b.Namespace), lorm.C("Name").EQ(b.Name)).Exec(context.Background()).Err()
}

func (r RoleBindingDao) FindByNameAndNamespace(name, namespace string) (rb *model.RoleBindings, err error){
	//err = r.db.Model(&model.RoleBindings{}).Where("name = ? and namespace = ?", name, namespace).Find(&rb).Error
	//return

	rb, err = lorm.NewSelector[model.RoleBindings](r.db).Where(lorm.C("namespace").EQ(namespace), lorm.C("Name").EQ(name)).Get(context.Background())
	return
}
