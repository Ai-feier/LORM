package dao

import (
	"context"
	"github.com/Ai-feier/lorm"
	"github.com/Ai-feier/rbacapp/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"strings"
)

/*
	lorm
 */

var db *lorm.DB

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
	db, err = lorm.Open("mysql", dsn)
	if err != nil {
		logrus.Error("数据库连接失败:", err)
		return
	}
}

func GetDB(ctx context.Context) *lorm.DB {
	return db
}

/*
	gorm
 */
//
//// 全局维护一个数据库链接, 其余链接从当前链接衍生
//var db2 *gorm.DB
//
//
//
//func InitDB() {
//	// 连接数据库
//	host := config.Conf.Mysql.Host
//	port := config.Conf.Mysql.Port
//	database := config.Conf.Mysql.Database
//	username := config.Conf.Mysql.Username
//	password := config.Conf.Mysql.Password
//	charset := config.Conf.Mysql.Charset
//	dsn := strings.Join([]string{username, ":", password, "@tcp(", host, ":", port, ")/", database, "?charset=",
//		charset, "&parseTime=true"}, "")
//
//	var err error
//	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
//	if err != nil {
//		logrus.Error("数据库连接失败:", err)
//		return
//	}
//}
//
//
//func GetDB(ctx context.Context) *gorm.DB {
//	if ctx == nil {
//		ctx = context.Background()
//	}
//	return db2.WithContext(ctx)
//}
