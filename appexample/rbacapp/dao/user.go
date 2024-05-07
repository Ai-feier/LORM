package dao

import (
	"context"
	"fmt"
	"github.com/Ai-feier/lorm"
	"github.com/Ai-feier/rbacapp/model"
	"github.com/Ai-feier/rbacapp/pkg/errs"
	"github.com/sirupsen/logrus"
	"time"
)

// 创建用户
// 删除用户
// 更新用户表
// 用户列表

type UserDao struct {
	db *lorm.DB

	sl lorm.Selector[model.Users]
}

func NewUserDao(ctx context.Context) *UserDao {
	return &UserDao{db: GetDB(ctx)}
}

func (u *UserDao) CreateUser(user *model.Users) (err error) {
	// 判断用户名是否存在
	//var cnt int64
	//
	//u.db.Model(&model.Users{}).Where("user_name = ?", user.UserName).Count(&cnt)
	//if cnt > 0 {
	//	err = errs.ExistUserName(user.UserName)
	//	logrus.Error(err.Error())
	//	return
	//}

	get, _ := lorm.NewSelector[model.Users](u.db).Where(lorm.C("Name").EQ(user.UserName)).Get(context.Background())
	if get != nil {
			err = errs.ExistUserName(user.UserName)
			logrus.Error(err.Error())
			return
	}

	//
	//// 更新时间
	now := time.Now().Unix()
	user.CreateTime = now
	user.UpdateTime = now

	//err = u.db.Model(&model.Users{}).Create(&user).Error
	//if err != nil {
	//	logrus.Error("用户创建失败: ", err)
	//}
	//return

	return lorm.NewInserter[model.Users](u.db).Values(user).Exec(context.Background()).Err()

}

func (u *UserDao) DeleteUser(user *model.Users) (err error) {
	//return u.db.Model(&model.Users{}).Where("id = ?", user.ID).Delete(&model.Users{}).Error

	return lorm.NewDeleter[model.Users](u.db).Where(lorm.C("id").EQ(user.ID)).Exec(context.Background()).Err()
}

func (u *UserDao) UpdateUser(user *model.Users) (err error) {
	user.UpdateTime = time.Now().Unix()
	//err = u.db.Model(&model.Users{}).Where("id = ?", user.ID).Updates(&user).Error
	//if err != nil {
	//	logrus.Error("用户更新失败: ", err)
	//}
	//return

	// 写法一
	//return lorm.NewUpdater[model.Users](u.db).Where(lorm.C("id").EQ(user.ID)).Set(lorm.Assign("UpdateTime", user.UpdateTime)).Exec(context.Background()).Err()

	// 写法二
	return lorm.NewUpdater[model.Users](u.db).Update(user).Set(lorm.C("UpdateTime")).Exec(context.Background()).Err()
}

func (u *UserDao) ListUser() (users []*model.Users, err error) {
	//err = u.db.Model(&model.Users{}).Find(&users).Error
	//return

	users, err = lorm.NewSelector[model.Users](u.db).GetMulti(context.Background())
	return
}

func (u *UserDao) FindByName(name string) (user *model.Users, err error) {
	//err = u.db.Model(&model.Users{}).Where("user_name = ?", name).Find(&user).Error
	//return
	return u.sl.Where(lorm.C("Name").EQ(name)).Get(context.Background())
}

func (u *UserDao) FindByNameAndPass(name, pass string) (user *model.Users, err error) {
	//var cnt int64
	//err = u.db.Model(&model.Users{}).Where("user_name = ? and password = ?", name, pass).Find(&user).Count(&cnt).Error
	//if cnt == 0 {
	//	return nil, errors.New("用户不存在")
	//}
	//return

	fmt.Println(u.sl.Where(lorm.C("Name").EQ(name), lorm.C("Password").EQ(pass)).Build())

	return u.sl.Where(lorm.C("UserName").EQ(name), lorm.C("Password").EQ(pass)).Get(context.Background())
}


