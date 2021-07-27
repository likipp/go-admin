package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	MysqlAdmin MySQL  `json:"mysqlAdmin"`
	RedisAdmin Redis  `json:"redisAdmin"`
	JWT        JWT    `json:"jwt"`
	Casbin     Casbin `json:"casbin"`
}

type MySQL struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Path     string `json:"path"`
	DBName   string `json:"dbname"`
	Config   string `json:"config"`
}

type Redis struct {
	Path string `json:"path"`
}

type JWT struct {
	SigningKey string `mapstructure:"signing-key" json:"signingKey" yaml:"signing-key"`
}

type Casbin struct {
	ModelPath string `json:"modelPath"`
}

var AdminConfig Config
var VTool *viper.Viper

func init() {
	v := viper.New()
	v.SetConfigName("settings")
	v.AddConfigPath("./config/")
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("配置文件读取错误: %s \n", err))
	}
	if err := v.Unmarshal(&AdminConfig); err != nil {
		fmt.Println(err)
	}
	VTool = v
}
