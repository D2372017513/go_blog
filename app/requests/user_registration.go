package requests

import (
	"fmt"
	"goblog/app/models/user"

	"github.com/thedevsaddam/govalidator"
)

// ValidateRegistrationForm 验证表单，返回 errs 长度等于零即通过
func ValidateRegistrationForm(data user.User) map[string][]string {
	// 表单规则
	rules := govalidator.MapData{
		"name":            []string{"required", "alpha_num", "between:3,20", "not_exist:users,name"},
		"email":           []string{"required", "min:4", "max:30", "email", "not_exist:users,email"},
		"password":        []string{"required", "min:6"},
		"passwordConfirm": []string{"required"},
	}

	// 定制错误信息
	messages := govalidator.MapData{
		"name": []string{
			"required:用户名为必填项",
			"alpha_num:格式错误，只允许数字和英文",
			"between:用户名长度需在 3~20 之间",
		},
		"email": []string{
			"required:Email 为必填项",
			"min:Email 长度需大于 4",
			"max:Email 长度需小于 30",
			"email:Email 格式不正确，请提供有效的邮箱地址",
		},
		"password": []string{
			"required:密码为必填项",
			"min:长度需大于 6",
		},
		"password_confirm": []string{
			"required:确认密码框为必填项",
		},
	}
	opts := govalidator.Options{
		Data:          &data,
		Rules:         rules,
		TagIdentifier: "vaild",
		Messages:      messages,
	}

	errs := govalidator.New(opts).ValidateStruct()
	fmt.Println(data)
	if data.Password != data.PasswordConfirm {
		errs["password_confirm"] = []string{"两次输入密码不符"}
	}

	return errs
}
