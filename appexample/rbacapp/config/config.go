package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

var Conf *Config

type Config struct {
	Mysql *Mysql `yaml:"mysql"`
}

type Mysql struct {
	DriveName string `yaml:"driveName"`
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	Database  string `yaml:"database"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	Charset   string `yaml:"charset"`
}

func InitConfig() {
	wordDir, _ := os.Getwd()
	viper.SetDefault("config_path", "./config")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	viper.AddConfigPath(wordDir+viper.GetString("config_path"))
	viper.ReadInConfig()

	err := viper.Unmarshal(&Conf)
	if err != nil {
		logrus.Panic("配置文件解析失败: ", err)
		return
	}
}

// 固定配置文件
//func InitConfig() {
//	viper.SetDefault("configPath", "./config/")  // 设置默认配置文件位置
//	viper.BindEnv("configPath", "config-path")  // 绑定环境变量中配置文件位置
//	//workDir, _ := os.Getwd()
//	viper.SetConfigName("config")
//	viper.SetConfigType("yaml")
//	viper.AddConfigPath(viper.GetString("configPath"))
//	err := viper.ReadInConfig()
//	if err != nil {
//		logrus.Panic("获取配置文件失败: ", err)
//		return
//	}
//
//	err = viper.Unmarshal(&Conf)
//	if err != nil {
//		logrus.Panic("配置文件解析失败: ", err)
//		return
//	}
//}

