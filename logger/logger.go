package logger

import "log"

// LogErr 当存在错误时记录日志
func LogErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
