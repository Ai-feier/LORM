package dao

import (
	"context"
	"github.com/Ai-feier/rbacapp/config"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
)

// 全局维护一个数据库链接, 其余链接从当前链接衍生
var db *gorm.DB



func InitDB() {
	// 连接数据库
	host := config.Conf.Mysql.Host
	port := config.Conf.Mysql.Port
	database := config.Conf.Mysql.Database
	username := config.Conf.Mysql.Username
	password := config.Conf.Mysql.Password
	charset := config.Conf.Mysql.Charset
	dsn := strings.Join([]string{username, ":", password, "@tcp(", host, ":", port, ")/", database, "?charset=",
		charset, "&parseTime=true"}, "")

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Error("数据库连接失败:", err)
		return
	}
}


func GetDB(ctx context.Context) *gorm.DB {
	if ctx == nil {
		ctx = context.Background()
	}
	return db.WithContext(ctx)
}
