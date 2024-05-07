package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/Ai-feier/lorm"
	"github.com/Ai-feier/rbacapp/model"
)

type RoleDao struct {
	db *lorm.DB
}

func NewRoleDao(ctx context.Context) *RoleDao {
	return &RoleDao{db: GetDB(ctx)}
}

func (r *RoleDao) FindByNameAndNameSpace(name, namespace string) (role *model.Roles, err error) {
	//var cnt int64
	//err = r.db.Model(&model.Roles{}).Where("name = ? and namespace = ?", name, namespace).Find(&role).Count(&cnt).Error
	//if cnt == 0 {
	//	return nil, errors.New(fmt.Sprintf("当前 %s 下不存在 %s role", namespace, name))
	//}

	//fmt.Println(lorm.NewSelector[model.Roles](r.db).Where(lorm.C("Name").EQ(name), lorm.C("Namespace").EQ(namespace)).Build())
	role, err = lorm.NewSelector[model.Roles](r.db).Where(lorm.C("Name").EQ(name), lorm.C("Namespace").EQ(namespace)).Get(context.Background())
	fmt.Println(role)
	if err != nil {
			return nil, errors.New(fmt.Sprintf("当前 %s 下不存在 %s role", namespace, name))
	}
	return
}

func (r *RoleDao) CreateRole(role *model.Roles) error {
	// 判断同一 namespace 下是否重名
	//var cnt int64
	//r.db.Model(&model.Roles{}).Where("name = ?", role.Name).Count(&cnt)
	//if cnt > 0 {
	//	return fmt.Errorf("当前命名空间 %s 下已存在名为 %s 的 role, 请重命名", role.Namespace, role.Name)
	//}
	//return r.db.Model(&model.Roles{}).Create(&role).Error
	res, _ := lorm.NewSelector[model.Roles](r.db).Where(lorm.C("Name").EQ(role.Name)).Get(context.Background())
	if res != nil {
			return fmt.Errorf("当前命名空间 %s 下已存在名为 %s 的 role, 请重命名", role.Namespace, role.Name)
	}

	ans := lorm.NewInserter[model.Roles](r.db).Values(role).Exec(context.Background())
	return ans.Err()
}

func (r *RoleDao) DeleteRole(role *model.Roles) error {
	//return r.db.Model(&model.Roles{}).Where("name = ? and namespace = ?", role.Name, role.Namespace).
		//Delete(&model.Roles{}).Error

	// 删除前检查是否存在
	//var cnt int64
	//r.db.Model(&model.Roles{}).Where("name = ? and namespace = ?", role.Name, role.Namespace).Count(&cnt)
	//if cnt == 0 {
	//	return fmt.Errorf("前命名空间 %s 下不存在名为 %s 的 role", role.Namespace, role.Name)
	//}
	//return r.db.Model(&model.Roles{}).Where("name = ? and namespace = ?", role.Name, role.Namespace).Delete(&model.Roles{}).Error

	fmt.Println(lorm.NewDeleter[model.Roles](r.db).Where(lorm.C("Name").EQ(role.Name), lorm.C("Namespace").EQ(role.Namespace)).Build())
	res := lorm.NewDeleter[model.Roles](r.db).Where(lorm.C("Name").EQ(role.Name), lorm.C("Namespace").EQ(role.Namespace)).Exec(context.Background())
	return res.Err()
}

func (r *RoleDao) DeleteRoleByNameAndNamespace(name, namespace string) error {
	//return r.db.Model(&model.Roles{}).Where("name = ? and namespace = ?", name, namespace).
	//	Delete(&model.Roles{}).Error

	res := lorm.NewDeleter[model.Roles](r.db).Where(lorm.C("Name").EQ(name), lorm.C("Namespace").EQ(namespace)).Exec(context.Background())
	return res.Err()
}
