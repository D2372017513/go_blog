package controllers

import (
	"fmt"
	"goblog/app/models/article"
	"goblog/app/models/user"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/view"
	"net/http"

	"gorm.io/gorm"
)

type UserController struct {
}

func (uc UserController) Show(w http.ResponseWriter, r *http.Request) {
	user_id := route.GetRouteVariable("id", r)

	user, err := user.Get(user_id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 用户未找到")
		} else {
			logger.LogErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		articles, err := article.GetByUserID(user_id)
		if err != nil {
			logger.LogErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		} else {
			view.Render(w, view.D{
				"Articles": articles,
				"User":     user,
			}, "articles.index", "articles._article_meta")
		}
	}

}
