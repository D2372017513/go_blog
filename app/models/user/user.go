package user

import (
	"goblog/app/models"
	"goblog/pkg/logger"
	"goblog/pkg/model"
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
