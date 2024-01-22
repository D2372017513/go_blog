package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

var v *viper.Viper = viper.New()
var cfg Config

type ServerCfg struct {
	Port int
}

type Config struct {
	ServerCfg   ServerCfg    `yaml:"server" mapstructure:"server"`
	DatabaseCfg mysql.Config `yaml:"database" mapstructure:"database"`
}

func init() {
	initCfg()
}

// 初始化配置文件
func initCfg() {
	addConfigPath()
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		// TODO 当配置文件发生变化之后进行处理
		fmt.Printf("there is something has changed!: %s\n", in.Name)
	})

	err = v.Unmarshal(&cfg)
	if err != nil {
		panic(fmt.Errorf("unmarshal error %v", err.Error()))
	}
}

func GetCfgByName(key string) any {
	return v.Get(key)
}

func GetDBCfg() mysql.Config {
	return cfg.DatabaseCfg
}

func GetServerCfg() ServerCfg {
	return cfg.ServerCfg
}

// 添加所有的路径
func addConfigPath() {
	v.SetConfigType("yaml")
	v.SetConfigFile("./config/config.yaml")
	v.AddConfigPath(".")
}
