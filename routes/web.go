package routes

import (
	"goblog/app/http/controllers"
	"goblog/app/http/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterWebRoutes 注册相关路由
func RegisterWebRoutes(router *mux.Router) {
	pc := new(controllers.PagesController)
	router.HandleFunc("/", pc.Home).Methods("GET").Name("home")
	router.HandleFunc("/about", pc.About).Methods("GET").Name("about")

	// 文章相关界面
	ac := new(controllers.ArticleController)
	router.HandleFunc("/articles/{id:[0-9]+}", ac.Show).Methods("GET").Name("articles.show")

	// 文章列表
	router.HandleFunc("/articles", ac.Index).Methods("GET").Name("articles.index")

	// 文章创建
	router.HandleFunc("/articles", ac.Store).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/create", ac.Create).Methods("GET").Name("articles.create")

	// 文章更新
	router.HandleFunc("/articles/{id:[0-9]+}/edit", ac.Edit).Methods("GET").Name("articles.edit")
	router.HandleFunc("/articles/{id:[0-9]+}", ac.Update).Methods("POST").Name("articles.update")

	// 文章删除
	router.HandleFunc("/articles/{id:[0-9]+}/delete", ac.Delete).Methods("POST").Name("articles.delete")

	// 自定义 404 页面
	router.NotFoundHandler = http.HandlerFunc(pc.NotFound)

	// 中间件：强制内容类型为 HTML
	router.Use(middlewares.ForceHTML)
}
