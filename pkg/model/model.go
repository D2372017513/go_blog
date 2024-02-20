package model

import (
	"goblog/config"
	"goblog/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB {
	var err error
	config := mysql.New(mysql.Config{
		DSN: config.GetMysqlPath(),
	})

	// 连接数据库
	DB, err = gorm.Open(config, &gorm.Config{})

	logger.LogErr(err)

	return DB
}
