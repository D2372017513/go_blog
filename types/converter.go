package types

import (
	"goblog/pkg/logger"
	"strconv"
)

// Int64ToString 将 int64 转换为 string
func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
}

// StringToUint64 将字符串转换为 uint64
func StringToUint64(str string) uint64 {
	i, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		logger.LogErr(err)
	}
	return i
}

// StringToInt64 将字符串转换为 int64
func StringToInt64(str string) int64 {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		logger.LogErr(err)
	}
	return i
}