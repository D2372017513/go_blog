package database

import (
	"database/sql"
	"fmt"
	"goblog/config"
	"goblog/logger"
	"log"
	"os"
	"strings"
	"time"
)

var db *sql.DB

// 初始化数据库
func init() {
	// 初始化连接数据库
	initDB()

	// 初始化数据库
	createDB()
}

func initDB() {
	var err error
	dbConfig := config.GetDBCfg()
	db, err = sql.Open("mysql", dbConfig.FormatDSN())
	logger.LogErr(err)

	// 设置最大连接数
	db.SetMaxOpenConns(25)
	// 设置最大空闲连接数
	db.SetMaxIdleConns(25)
	// 设置每个链接的过期时间
	db.SetConnMaxLifetime(5 * time.Minute)

	// 尝试连接，失败会报错
	err = db.Ping()
	logger.LogErr(err)
}

func createDB() error {
	workPath, _ := os.Getwd()
	sqlPath := workPath + "\\database\\database.sql"
	if !checkFileExists(sqlPath) {
		return fmt.Errorf("sql file is not exist")
	}

	sqls, _ := os.ReadFile(sqlPath)
	sqlArr := strings.Split(string(sqls), ";")
	for _, sql := range sqlArr {
		sql = strings.TrimSpace(sql)
		if sql == "" {
			continue
		}

		result, err := db.Exec(sql)
		if err != nil {
			log.Println("数据库导入失败:" + err.Error())
			return err
		} else {
			lastId, _ := result.LastInsertId()
			affectRow, _ := result.RowsAffected()
			log.Printf("\t %s Exec success! lastId : %d, affectRow: %d", sql, lastId, affectRow)
		}
	}
	return nil
}

// 返回数据库连接对象
func GetDB() *sql.DB {
	return db
}

func checkFileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
