package requests

import (
	"errors"
	"fmt"
	"goblog/pkg/model"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/thedevsaddam/govalidator"
)

func init() {
	govalidator.AddCustomRule("not_exist", func(field, rule, message string, value interface{}) error {
		rng := strings.Split(strings.TrimPrefix(rule, "not_exist:"), ",")
		tableName, columnName := rng[0], rng[1]

		var count int64
		model.DB.Table(tableName).Where(columnName+" = ?", value.(string)).Count(&count)

		if count != 0 {
			if message != "" {
				return errors.New(message)
			}

			return fmt.Errorf("%v 已被占用", value)
		}

		return nil
	})

	// max_cn:8
	govalidator.AddCustomRule("max_cn", func(field string, rule string, message string, value interface{}) error {
		valLength := utf8.RuneCountInString(value.(string))
		l, _ := strconv.Atoi(strings.TrimPrefix(rule, "max_cn:")) //handle other error
		if valLength > l {
			if message != "" {
				return errors.New(message)
			}
			return fmt.Errorf("长度不能超过 %d 个字", l)
		}
		return nil
	})

	// min_cn:2
	govalidator.AddCustomRule("min_cn", func(field string, rule string, message string, value interface{}) error {
		valLength := utf8.RuneCountInString(value.(string))
		l, _ := strconv.Atoi(strings.TrimPrefix(rule, "min_cn:")) //handle other error
		if valLength < l {
			if message != "" {
				return errors.New(message)
			}
			return fmt.Errorf("长度需大于 %d 个字", l)
		}
		return nil
	})
}
