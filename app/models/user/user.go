package user

import (
	"goblog/app/models"
	"goblog/pkg/logger"
	"goblog/pkg/model"
	"goblog/pkg/password"
)

type User struct {
	models.BaseModel

	Name     string `gorm:"column:name;type:varchar(255);not null;unique" vaild:"name"`
	Email    string `gorm:"column:email;type:varchar(255);default:null;unique" vaild:"email"`
	Password string `gorm:"column:password;type:varchar(255)" vaild:"password"`

	PasswordConfirm string `gorm:"-" vaild:"passwordConfirm"`
}

// Create 插入一个新用户
func (user *User) Create() (err error) {
	if err = model.DB.Create(&user).Error; err != nil {
		logger.LogErr(err)
		return err
	}

	return nil
}

// ComparePassword 比较用户密码
func (user *User) ComparePassword(passwd string) bool {
	return password.CheckHash(passwd, user.Password)
}

// CompareEmail 比较用户邮箱
func (u *User) CompareEmail(email string) bool {
	return email == u.Email
}
