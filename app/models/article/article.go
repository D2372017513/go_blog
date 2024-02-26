package article

import (
	"strconv"

	"goblog/app/models"
	"goblog/app/models/user"
	"goblog/pkg/route"
)

type ArticlesData struct {
	models.BaseModel
	Title      string `gorm:"type:varchar(255);not null" valid:"title"`
	Body       string `gorm:"type:longtext;not null" valid:"body"`
	UserID     int64  `gorm:"index"`
	User       user.User
	CategoryID uint64 `gorm:"not null;default:4;index"`
	URL        string `gorm:"-"`
}

// 修改 gorm 的默认表名
func (ArticlesData) TableName() string {
	return "articles"
}

func (a ArticlesData) Link() string {
	return route.Name2URL("articles.show", "id", strconv.FormatInt(a.ID, 10))
}

// CreatedAtDate 创建日期
func (article ArticlesData) CreatedAtDate() string {
	return article.CreateAt.Format("2006-01-02")
}
