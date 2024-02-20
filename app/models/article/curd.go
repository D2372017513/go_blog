package article

import (
	"goblog/logger"
	"goblog/pkg/model"
)

// 获取全部文章
func GetAll() ([]ArticlesData, error) {
	var articles []ArticlesData
	if err := model.DB.Find(&articles); err != nil {
		return articles, err.Error
	}

	return articles, nil
}

// Create 创建文章，通过 article.ID 来判断是否创建成功
func (article *ArticlesData) Create() (err error) {
	result := model.DB.Create(&article)
	if err = result.Error; err != nil {
		logger.LogErr(err)
		return err
	}

	return nil
}
