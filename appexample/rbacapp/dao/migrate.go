package dao

import (
	"github.com/Ai-feier/rbacapp/model"
	"github.com/sirupsen/logrus"
)

func Migrate() {
	err := db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.RoleBinding{},
		&model.RoleSubRef{},
		&model.ClusterRole{},
		&model.ClusterRoleBinding{},
		&model.ClusterRoleSubRef{},
	)
	if err != nil {
		logrus.Error("自动迁移表结构失败:", err)
		return
	}

	logrus.Info("表结构迁移成功！")
}
