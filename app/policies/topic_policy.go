package policies

import (
	"goblog/app/models/article"
	"goblog/pkg/auth"
)

// CanModifyArticle 是否允许修改话题
func CanModifyArticle(a article.ArticlesData) bool {
	return auth.User().ID == a.UserID
}
