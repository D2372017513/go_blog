package requests

import (
	"errors"
	"fmt"
	"goblog/pkg/model"
	"strings"

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
}
