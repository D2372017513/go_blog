package user

import (
	"goblog/pkg/model"
	"goblog/pkg/types"
)

// All 获取所有用户数据
func All() ([]User, error) {
	var users []User
	if err := model.DB.Find(&users).Error; err != nil {
		return users, err
	}
	return users, nil
}

func Get(uid string) (User, error) {
	var user User
	uidI := types.StringToInt64(uid)

	if err := model.DB.First(&user, uidI).Error; err != nil {
		return user, err
	}

	return user, nil
}

func GetByEmail(email string) (User, error) {
	var user User
	if err := model.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

// GetByName 通过用户名获取用户
func GetByName(name string) (User, error) {
	var user User
	if err := model.DB.Where("name = ?", name).First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}
