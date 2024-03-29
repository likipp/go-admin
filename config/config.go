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
	Logs       Logs   `json:"logs"`
}

type MySQL struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Path     string `json:"path"`
	DBName   string `json:"dbname"`
	Config   string `json:"config"`
}

type Redis struct {
	Path     string `json:"path"`
	Password string `json:"password"`
}

type JWT struct {
	SigningKey string `mapstructure:"signing-key" json:"signingKey" yaml:"signing-key"`
}

type Casbin struct {
	ModelPath string `json:"modelPath"`
}

type Logs struct {
	Path string `json:"path"`
	Name string `json:"name"`
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
