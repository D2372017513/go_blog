package article

import (
	"fmt"
	"goblog/pkg/logger"
	"goblog/pkg/model"
	"goblog/types"
)

func Get(idstr string) (ArticlesData, error) {
	article := ArticlesData{}
	id := types.StringToUint64(idstr)
	if err := model.DB.Preload("User").First(&article, id).Error; err != nil {
		return article, err
	}

	return article, nil
}

// 获取全部文章
func GetAll() ([]ArticlesData, error) {
	var articles []ArticlesData
	if err := model.DB.Preload("User").Find(&articles); err != nil {
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

// Update 创建文章
func (article *ArticlesData) Update() (rowsAffected int64, err error) {
	fmt.Println(article)
	result := model.DB.Save(&article)
	fmt.Println(article)
	if err = result.Error; err != nil {
		logger.LogErr(err)
		return 0, err
	}

	return result.RowsAffected, nil
}

// Delete 删除文章
func (article *ArticlesData) Delete() (rowsAffected int64, err error) {
	result := model.DB.Delete(&article)
	if err = result.Error; err != nil {
		logger.LogErr(err)
		return 0, err
	}

	return result.RowsAffected, nil
}
