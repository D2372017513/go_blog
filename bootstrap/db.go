package bootstrap

import (
	"database/sql"
	"fmt"
	"goblog/pkg/model"
	"log"
	"os"
	"strings"
	"time"
)

var sqlDB *sql.DB

func init() {
	SetupDB()
}

func SetupDB() {
	db := model.ConnectDB()

	sqlDB, _ = db.DB()

	// 设置最大连接数
	sqlDB.SetMaxOpenConns(100)
	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(25)
	// 设置每个链接的过期时间
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	createDB(sqlDB)
}

func createDB(db *sql.DB) error {
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
			log.Printf("\t %s Exec success! lastId : %d, affectRow: %d", "sss", lastId, affectRow)
		}
	}
	return nil
}

func checkFileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
