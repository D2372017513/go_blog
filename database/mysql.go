package database

import (
	"database/sql"
	"goblog/bootstrap"
)

// 返回数据库连接对象
func GetDB() *sql.DB {
	return bootstrap.GetDB()
}
