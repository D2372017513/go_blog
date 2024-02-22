package user

import (
	"goblog/pkg/model"
	"goblog/types"
)

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
