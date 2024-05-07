package service

import (
	"context"
	"fmt"
	"github.com/Ai-feier/rbacapp/dao"
	"github.com/Ai-feier/rbacapp/model"
)

type UserSvc struct {
}

func NewUserSvc() *UserSvc {
	return &UserSvc{}
}

func (s *UserSvc) Login(ctx context.Context, username, password string) (user *model.Users, err error){
	userDao := dao.NewUserDao(ctx)
	user, err = userDao.FindByNameAndPass(username, password)
	if err != nil {
		return nil, fmt.Errorf("username: %s 不存在: %w", username, err)
	}
	return
}

func (s *UserSvc) Registry(ctx context.Context, username, password string) (user *model.Users, err error) {
	userDao := dao.NewUserDao(ctx)
	//user, err = userDao.FindByName(username)
	//if user != nil {
	//	return nil, fmt.Errorf("username: %s 已存在", username)
	//}
	//if err != nil {
	//	return nil, err
	//}
	user = &model.Users{}
	user.UserName = username
	user.Password = password
	err = userDao.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("用户注册失败: %w", err)
	}
	return
}
