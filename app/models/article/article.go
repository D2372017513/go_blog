package article

import (
	"strconv"

	"goblog/app/models"
	"goblog/pkg/model"
	"goblog/pkg/route"
	"goblog/types"
)

type ArticlesData struct {
	models.BaseModel
	Title string `gorm:"type:varchar(255);not null" valid:"title"`
	Body  string `gorm:"type:longtext;not null" valid:"body"`
	URL   string `gorm:"-"`
}

// 修改 gorm 的默认表名
func (ArticlesData) TableName() string {
	return "articles"
}

func (a *ArticlesData) Link() string {
	return route.Name2URL("articles.show", "id", strconv.FormatInt(a.ID, 10))
}

func Get(idstr string) (ArticlesData, error) {
	article := ArticlesData{}
	id := types.StringToUint64(idstr)
	if err := model.DB.First(&article, id).Error; err != nil {
		return article, err
	}

	return article, nil
}
