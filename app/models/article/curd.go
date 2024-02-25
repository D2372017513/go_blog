package article

import (
	"fmt"
	"goblog/pkg/logger"
	"goblog/pkg/model"
	"goblog/pkg/pagination"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"net/http"
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
func GetAll(r *http.Request, perPage int) ([]ArticlesData, pagination.ViewData, error) {
	// 1. 初始化分页实例
	db := model.DB.Model(ArticlesData{}).Order("create_at desc")
	_pager := pagination.New(r, db, route.Name2URL("home"), perPage)

	// 2. 获取视图数据
	viewData := _pager.Paging()

	// 3. 获取数据
	var articles []ArticlesData
	_pager.Results(&articles)

	return articles, viewData, nil
}

// GetByUserID 通过用户id获取文章
func GetByUserID(user_id string) (articles []ArticlesData, err error) {
	if err = model.DB.Where("user_id = ?", user_id).Preload("User").Find(&articles).Error; err != nil {
		return
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
