package bootstrap

import (
	"database/sql"
	"goblog/app/models/article"
	"goblog/app/models/category"
	"goblog/app/models/user"
	"goblog/pkg/config"
	"goblog/pkg/model"
	"time"

	"gorm.io/gorm"
)

var sqlDB *sql.DB

func init() {
	SetupDB()
}

func SetupDB() {
	db := model.ConnectDB()

	sqlDB, _ = db.DB()

	// 设置最大连接数
	sqlDB.SetMaxOpenConns(config.GetInt("database.mysql.max_open_connections"))
	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(config.GetInt("database.mysql.max_idle_connections"))
	// 设置每个链接的过期时间
	sqlDB.SetConnMaxLifetime(time.Duration(config.GetInt("database.mysql.max_life_seconds")) * time.Second)

	// 创建和维护数据表结构
	migration(db)
}

func migration(db *gorm.DB) {
	db.AutoMigrate(
		&article.ArticlesData{},
		&user.User{},
		&category.Category{},
	)
}
