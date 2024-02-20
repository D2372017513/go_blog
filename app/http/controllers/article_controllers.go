package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"goblog/app/models/article"
	"goblog/logger"
	"goblog/pkg/route"
	"goblog/types"

	"gorm.io/gorm"
)

type ArticlesData struct {
	ID          int64
	Title, Body string
	URL         *url.URL
	// Errors      map[string]string
	ShowErr string
}

type ArticleController struct {
}

func (ac *ArticleController) Show(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	rs, err := article.Get(id)

	// 3. 如果出现错误
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			// 3.2 数据库错误
			logger.LogErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		// 4. 读取成功
		tmpl, err := template.New("show.gohtml").Funcs(template.FuncMap{
			"RouteName2URL": route.Name2URL,
			"Int64ToString": types.Int64ToString,
		}).ParseFiles("resources/views/articles/show.gohtml")
		logger.LogErr(err)
		err = tmpl.Execute(w, rs)
		logger.LogErr(err)
	}
}
