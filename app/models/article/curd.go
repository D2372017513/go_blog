package article

import "goblog/pkg/model"

// 获取全部文章
func GetAll() ([]ArticlesData, error) {
	var articles []ArticlesData
	if err := model.DB.Find(&articles); err != nil {
		return articles, err.Error
	}

	return articles, nil
}
