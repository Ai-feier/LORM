package dao

import (
	"context"
	"errors"
	"github.com/Ai-feier/rbacapp/model"
	"github.com/Ai-feier/rbacapp/pkg/errs"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

// 创建用户
// 删除用户
// 更新用户表
// 用户列表

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(ctx context.Context) *UserDao {
	return &UserDao{db: GetDB(ctx)}
}

func (u *UserDao) CreateUser(user *model.User) (err error) {
	// 判断用户名是否存在
	var cnt int64

	u.db.Model(&model.User{}).Where("user_name = ?", user.UserName).Count(&cnt)
	if cnt > 0 {
		err = errs.ExistUserName(user.UserName)
		logrus.Error(err.Error())
		return
	}

	// 更新时间
	now := time.Now().Unix()
	user.CreateTime = now
	user.UpdateTime = now
	err = u.db.Model(&model.User{}).Create(&user).Error
	if err != nil {
		logrus.Error("用户创建失败: ", err)
	}
	return
}

func (u *UserDao) DeleteUser(user *model.User) (err error) {
	return u.db.Model(&model.User{}).Where("id = ?", user.ID).Delete(&model.User{}).Error
}

func (u *UserDao) UpdateUser(user *model.User) (err error) {
	user.UpdateTime = time.Now().Unix()
	err = u.db.Model(&model.User{}).Where("id = ?", user.ID).Updates(&user).Error
	if err != nil {
		logrus.Error("用户更新失败: ", err)
	}
	return
}

func (u *UserDao) ListUser() (users []*model.User, err error) {
	err = u.db.Model(&model.User{}).Find(&users).Error
	return
}

func (u *UserDao) FindByName(name string) (user *model.User, err error) {
	err = u.db.Model(&model.User{}).Where("user_name = ?", name).Find(&user).Error
	return
}

func (u *UserDao) FindByNameAndPass(name, pass string) (user *model.User, err error) {
	var cnt int64
	err = u.db.Model(&model.User{}).Where("user_name = ? and password = ?", name, pass).Find(&user).Count(&cnt).Error
	if cnt == 0 {
		return nil, errors.New("用户不存在")
	}
	return
}


